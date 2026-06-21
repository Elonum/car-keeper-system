import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/lib/AuthContext';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import EmptyState from '../common/EmptyState';
import CabinetListToolbar from './CabinetListToolbar';
import {
  FileText,
  Upload,
  Trash2,
  Download,
  AlertCircle,
  Paperclip,
  Link2,
  UserRound,
} from 'lucide-react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { getApiErrorMessage } from '@/lib/apiErrors';
import {
  documentService,
  DOCUMENT_TYPE_LABELS,
} from '@/services/documentService';
import { queryKeys } from '@/lib/queryKeys';
import {
  DOCUMENT_ACCEPT,
  SUPPORTED_FORMATS_LABEL,
  validateDocumentFile,
  formatFileSize,
} from '@/lib/documents';
import { ErrorNotice } from '@/components/common/ErrorNotice';
import { PERMISSIONS, hasPermission } from '@/lib/authz';
import ConfirmDialog from '@/components/common/ConfirmDialog';

const DOC_TYPES = Object.keys(DOCUMENT_TYPE_LABELS);

function docId(doc) {
  return doc.document_id?.toString?.() || doc.document_id;
}

function shortId(value) {
  return String(value || '').slice(0, 8);
}

function ownerLabelFromDoc(doc) {
  return doc.owner_name || doc.owner_email || 'Клиент';
}

function attachmentLabelFromDoc(doc) {
  if (doc.attachment_label) return doc.attachment_label;
  if (doc.order_id) return `Заказ #${shortId(doc.order_id)}`;
  if (doc.service_appointment_id) return `Запись на ТО #${shortId(doc.service_appointment_id)}`;
  return 'Привязка не указана';
}

function orderOptionLabel(order) {
  const id = order.order_id?.toString?.() || order.order_id;
  const customer = order.customer_name || order.customer_email;
  const status = order.status_label || order.status;
  return [
    `Заказ #${shortId(id)}`,
    customer,
    status,
  ].filter(Boolean).join(' · ');
}

function appointmentOptionLabel(appointment) {
  const id =
    appointment.service_appointment_id?.toString?.() ||
    appointment.service_appointment_id;
  const customer = appointment.owner_name || appointment.owner_email;
  const vin = appointment.user_car_vin ? `VIN ${appointment.user_car_vin}` : '';
  return [
    `Запись #${shortId(id)}`,
    customer,
    appointment.branch_name,
    vin,
  ].filter(Boolean).join(' · ');
}

export default function MyDocuments({ orders = [], appointments = [] }) {
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const fileInputRef = useRef(null);
  const [file, setFile] = useState(null);
  const [dragOver, setDragOver] = useState(false);
  const [documentType, setDocumentType] = useState('order_contract');
  const [linkKind, setLinkKind] = useState('order');
  const [orderId, setOrderId] = useState('');
  const [appointmentId, setAppointmentId] = useState('');
  const [formError, setFormError] = useState(null);
  const [actionError, setActionError] = useState(null);
  const [search, setSearch] = useState('');
  const [typeFilter, setTypeFilter] = useState('all');
  const [downloadingId, setDownloadingId] = useState(null);
  const [documentToDelete, setDocumentToDelete] = useState(null);

  const canLinkOrder = orders.length > 0;
  const canLinkAppointment = appointments.length > 0;
  const canUpload = canLinkOrder || canLinkAppointment;
  const staffDocumentsView = hasPermission(user, PERMISSIONS.DOCUMENTS_VIEW_ANY);

  useEffect(() => {
    if (!canLinkOrder && canLinkAppointment) setLinkKind('appointment');
    if (canLinkOrder && !canLinkAppointment) setLinkKind('order');
  }, [canLinkOrder, canLinkAppointment]);

  const { data: documents, isLoading } = useQuery({
    queryKey: queryKeys.myDocuments(),
    queryFn: () => documentService.list(),
  });

  const uploadMutation = useMutation({
    mutationFn: (payload) => documentService.upload(payload),
    onSuccess: (createdDocument) => {
      queryClient.setQueryData(queryKeys.myDocuments(), (current) => {
        const existing = Array.isArray(current) ? current : [];
        const createdId = docId(createdDocument);
        if (!createdId) return existing;
        return [
          createdDocument,
          ...existing.filter((doc) => docId(doc) !== createdId),
        ];
      });
      queryClient.invalidateQueries({ queryKey: queryKeys.myDocuments() });
      setFile(null);
      if (fileInputRef.current) fileInputRef.current.value = '';
      setFormError(null);
    },
    onError: (e) => {
      setFormError(getApiErrorMessage(e, 'Не удалось загрузить файл'));
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id) => documentService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.myDocuments() });
      setActionError(null);
    },
    onError: (e) => {
      setActionError(getApiErrorMessage(e, 'Не удалось удалить документ'));
    },
  });

  const pickFile = (next) => {
    if (!next) {
      setFile(null);
      return;
    }
    const check = validateDocumentFile(next);
    if (!check.ok) {
      setFormError(check.message);
      setFile(null);
      if (fileInputRef.current) fileInputRef.current.value = '';
      return;
    }
    setFormError(null);
    setFile(next);
  };

  const handleUpload = (e) => {
    e.preventDefault();
    const check = validateDocumentFile(file);
    if (!check.ok) {
      setFormError(check.message);
      return;
    }
    const payload = { file, documentType };
    if (linkKind === 'order') {
      if (!orderId) {
        setFormError('Выберите заказ');
        return;
      }
      payload.orderId = orderId;
    } else {
      if (!appointmentId) {
        setFormError('Выберите запись на ТО');
        return;
      }
      payload.serviceAppointmentId = appointmentId;
    }
    setFormError(null);
    uploadMutation.mutate(payload);
  };

  const list = Array.isArray(documents) ? documents : [];

  const filterOptions = useMemo(() => {
    const head = [{ value: 'all', label: 'Все документы' }];
    const rest = DOC_TYPES.map((t) => ({
      value: `type:${t}`,
      label: DOCUMENT_TYPE_LABELS[t],
    }));
    return [
      ...head,
      ...rest,
      { value: 'attach:order', label: 'Прикреплены к заказам' },
      { value: 'attach:service_appointment', label: 'Прикреплены к ТО' },
      { value: 'availability:missing', label: 'Файлы отсутствуют' },
    ];
  }, []);

  const filteredList = useMemo(() => {
    const q = search.trim().toLowerCase();
    return list.filter((doc) => {
      if (typeFilter.startsWith('type:') && doc.document_type !== typeFilter.slice(5)) {
        return false;
      }
      if (
        typeFilter.startsWith('attach:') &&
        (doc.attachment_kind || '') !== typeFilter.slice(7)
      ) {
        return false;
      }
      if (typeFilter === 'availability:missing' && doc.file_available !== false) {
        return false;
      }
      if (!q) return true;
      const name = (doc.file_name || DOCUMENT_TYPE_LABELS[doc.document_type] || '').toLowerCase();
      const typeLabel = (DOCUMENT_TYPE_LABELS[doc.document_type] || doc.document_type || '').toLowerCase();
      const id = String(docId(doc) || '').toLowerCase();
      const owner = `${doc.owner_name || ''} ${doc.owner_email || ''}`.toLowerCase();
      const attachment = attachmentLabelFromDoc(doc).toLowerCase();
      return (
        name.includes(q) ||
        typeLabel.includes(q) ||
        id.includes(q) ||
        owner.includes(q) ||
        attachment.includes(q)
      );
    });
  }, [list, search, typeFilter]);

  const stats = useMemo(() => {
    const total = list.length;
    const available = list.filter((doc) => doc.file_available !== false).length;
    const ordersCount = list.filter((doc) => doc.attachment_kind === 'order' || doc.order_id).length;
    const appointmentsCount = list.filter(
      (doc) => doc.attachment_kind === 'service_appointment' || doc.service_appointment_id
    ).length;
    return { total, available, ordersCount, appointmentsCount };
  }, [list]);

  const handleDownload = async (doc) => {
    const id = docId(doc);
    const name = doc.file_name || DOCUMENT_TYPE_LABELS[doc.document_type] || 'document';
    setDownloadingId(id);
    setActionError(null);
    try {
      await documentService.download(id, name);
    } catch (err) {
      setActionError(getApiErrorMessage(err, 'Не удалось скачать файл'));
    } finally {
      setDownloadingId(null);
    }
  };

  const confirmDeleteDocument = () => {
    if (!documentToDelete) return;
    deleteMutation.mutate(docId(documentToDelete), {
      onSettled: () => setDocumentToDelete(null),
    });
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2].map((i) => (
          <Card key={i} className="p-6 animate-pulse h-24 rounded-2xl border-slate-200" />
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Card className="rounded-2xl border-slate-200 p-6 shadow-sm">
        <div className="flex items-start gap-3 mb-5">
          <div className="w-10 h-10 rounded-xl bg-indigo-50 flex items-center justify-center flex-shrink-0">
            <Upload className="w-5 h-5 text-indigo-600" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-slate-900">Загрузить документ</h3>
            <p className="text-sm text-slate-500 mt-0.5">
              {SUPPORTED_FORMATS_LABEL} · до {formatFileSize(15 * 1024 * 1024)}
            </p>
          </div>
        </div>

        {!canUpload ? (
          <div className="flex items-start gap-2 rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900">
            <AlertCircle className="w-4 h-4 mt-0.5 flex-shrink-0" />
            <p>
              Чтобы загрузить файл, нужен хотя бы один заказ или запись на ТО — документ привязывается к
              ним.
            </p>
          </div>
        ) : (
          <form onSubmit={handleUpload} className="space-y-5 max-w-2xl">
            <ErrorNotice kind="form" message={formError} />

            <div className="space-y-2">
              <Label>Файл</Label>
              <div
                role="button"
                tabIndex={0}
                onKeyDown={(ev) => {
                  if (ev.key === 'Enter' || ev.key === ' ') fileInputRef.current?.click();
                }}
                onDragOver={(ev) => {
                  ev.preventDefault();
                  setDragOver(true);
                }}
                onDragLeave={() => setDragOver(false)}
                onDrop={(ev) => {
                  ev.preventDefault();
                  setDragOver(false);
                  pickFile(ev.dataTransfer.files?.[0] || null);
                }}
                onClick={() => fileInputRef.current?.click()}
                className={`cursor-pointer rounded-xl border-2 border-dashed px-6 py-8 text-center transition-colors ${
                  dragOver
                    ? 'border-indigo-400 bg-indigo-50/50'
                    : 'border-slate-200 bg-slate-50/80 hover:border-slate-300 hover:bg-slate-50'
                }`}
              >
                <Paperclip className="w-8 h-8 text-slate-400 mx-auto mb-2" />
                {file ? (
                  <p className="font-medium text-slate-900 truncate">{file.name}</p>
                ) : (
                  <p className="font-medium text-slate-700">Перетащите файл или нажмите для выбора</p>
                )}
                {file ? (
                  <p className="text-sm text-slate-500 mt-1">{formatFileSize(file.size)}</p>
                ) : null}
              </div>
              <input
                ref={fileInputRef}
                type="file"
                accept={DOCUMENT_ACCEPT}
                className="sr-only"
                onChange={(ev) => pickFile(ev.target.files?.[0] || null)}
              />
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label>Тип документа</Label>
                <Select
                  value={documentType}
                  onValueChange={(v) => {
                    setDocumentType(v);
                    setFormError(null);
                  }}
                >
                  <SelectTrigger className="bg-white">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {DOC_TYPES.map((t) => (
                      <SelectItem key={t} value={t}>
                        {DOCUMENT_TYPE_LABELS[t]}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Привязка</Label>
                <Select
                  value={linkKind}
                  onValueChange={(v) => {
                    setLinkKind(v);
                    setFormError(null);
                  }}
                >
                  <SelectTrigger className="bg-white">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {canLinkOrder ? <SelectItem value="order">К заказу</SelectItem> : null}
                    {canLinkAppointment ? (
                      <SelectItem value="appointment">К записи на ТО</SelectItem>
                    ) : null}
                  </SelectContent>
                </Select>
              </div>
            </div>

            {linkKind === 'order' && canLinkOrder ? (
              <div className="space-y-2">
                <Label>Заказ</Label>
                <Select
                  value={orderId}
                  onValueChange={(v) => {
                    setOrderId(v);
                    setFormError(null);
                  }}
                >
                  <SelectTrigger className="bg-white">
                    <SelectValue placeholder="Выберите заказ" />
                  </SelectTrigger>
                  <SelectContent>
                    {orders.map((o) => {
                      const id = o.order_id?.toString?.() || o.order_id;
                      return (
                        <SelectItem key={id} value={id}>
                          {orderOptionLabel(o)}
                        </SelectItem>
                      );
                    })}
                  </SelectContent>
                </Select>
              </div>
            ) : null}

            {linkKind === 'appointment' && canLinkAppointment ? (
              <div className="space-y-2">
                <Label>Запись на ТО</Label>
                <Select
                  value={appointmentId}
                  onValueChange={(v) => {
                    setAppointmentId(v);
                    setFormError(null);
                  }}
                >
                  <SelectTrigger className="bg-white">
                    <SelectValue placeholder="Выберите запись" />
                  </SelectTrigger>
                  <SelectContent>
                    {appointments.map((a) => {
                      const id =
                        a.service_appointment_id?.toString?.() || a.service_appointment_id;
                      return (
                        <SelectItem key={id} value={id}>
                          {appointmentOptionLabel(a)}
                        </SelectItem>
                      );
                    })}
                  </SelectContent>
                </Select>
              </div>
            ) : null}

            <Button type="submit" disabled={uploadMutation.isPending || !file}>
              {uploadMutation.isPending ? 'Загрузка…' : 'Загрузить'}
            </Button>
          </form>
        )}
      </Card>

      <div className="space-y-4">
        <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h3 className="text-lg font-semibold text-slate-900">Мои документы</h3>
            <p className="text-sm text-slate-500">
              {staffDocumentsView
                ? 'Показаны документы клиентов, доступные вашей роли.'
                : 'Документы, прикреплённые к вашим заказам и записям на ТО.'}
            </p>
          </div>
          <div className="flex flex-wrap gap-2 text-xs text-slate-600">
            <span className="rounded-full bg-slate-100 px-3 py-1">Всего: {stats.total}</span>
            <span className="rounded-full bg-emerald-50 text-emerald-700 px-3 py-1">
              Доступно: {stats.available}
            </span>
            <span className="rounded-full bg-blue-50 text-blue-700 px-3 py-1">
              Заказы: {stats.ordersCount}
            </span>
            <span className="rounded-full bg-violet-50 text-violet-700 px-3 py-1">
              ТО: {stats.appointmentsCount}
            </span>
          </div>
        </div>
        <ErrorNotice kind="server" message={actionError} />

        {list.length === 0 ? (
          <EmptyState
            icon={FileText}
            title="Нет документов"
            description="Загрузите файл, привязанный к заказу или записи на сервис"
          />
        ) : (
          <>
            <CabinetListToolbar
              search={search}
              onSearchChange={setSearch}
              searchPlaceholder="Поиск: файл, клиент, email, заказ, VIN…"
              statusValue={typeFilter}
              onStatusChange={setTypeFilter}
              statusOptions={filterOptions}
            />
            {filteredList.length === 0 ? (
              <p className="text-sm text-slate-500 py-4">Ничего не найдено по фильтрам.</p>
            ) : (
              <div className="space-y-3">
                {filteredList.map((doc) => {
                  const id = docId(doc);
                  const name =
                    doc.file_name || DOCUMENT_TYPE_LABELS[doc.document_type] || 'Файл';
                  const available = Boolean(doc.file_available);
                  const ownerLabel = ownerLabelFromDoc(doc);
                  const attachmentLabel = attachmentLabelFromDoc(doc);
                  const created =
                    doc.created_at &&
                    format(new Date(doc.created_at), 'd MMM yyyy, HH:mm', { locale: ru });

                  return (
                    <Card
                      key={id}
                      className="rounded-2xl border-slate-200 p-4 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4"
                    >
                      <div className="flex items-start gap-3 min-w-0 flex-1">
                        <div className="w-11 h-11 rounded-xl bg-slate-100 flex items-center justify-center flex-shrink-0">
                          <FileText className="w-5 h-5 text-slate-600" />
                        </div>
                        <div className="min-w-0">
                          <p className="font-medium text-slate-900 truncate">{name}</p>
                          <p className="text-sm text-slate-500">
                            {DOCUMENT_TYPE_LABELS[doc.document_type] || doc.document_type}
                            {doc.file_size != null ? ` · ${formatFileSize(doc.file_size)}` : ''}
                          </p>
                          <div className="mt-2 flex flex-col gap-1 text-xs text-slate-500 sm:flex-row sm:flex-wrap sm:gap-x-4">
                            <span className="inline-flex items-center gap-1 min-w-0">
                              <UserRound className="h-3.5 w-3.5 flex-shrink-0" />
                              <span className="truncate">
                                {ownerLabel}
                                {staffDocumentsView && doc.owner_email ? ` · ${doc.owner_email}` : ''}
                              </span>
                            </span>
                            <span className="inline-flex items-center gap-1 min-w-0">
                              <Link2 className="h-3.5 w-3.5 flex-shrink-0" />
                              <span className="truncate">{attachmentLabel}</span>
                            </span>
                          </div>
                          {created ? (
                            <p className="text-xs text-slate-400 mt-0.5">{created}</p>
                          ) : null}
                          {!available ? (
                            <p className="text-xs text-amber-700 mt-1">
                              Файл на сервере отсутствует (только запись в каталоге)
                            </p>
                          ) : null}
                        </div>
                      </div>
                      <div className="flex gap-2 w-full sm:w-auto sm:shrink-0 justify-end">
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          disabled={!available || downloadingId === id}
                          title={
                            available
                              ? 'Скачать'
                              : 'Файл недоступен — загрузите новый документ с тем же типом'
                          }
                          onClick={() => handleDownload(doc)}
                        >
                          <Download className="w-4 h-4 mr-1" />
                          {downloadingId === id ? 'Загрузка…' : 'Скачать'}
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          className="text-red-600 hover:text-red-700 hover:bg-red-50"
                          disabled={deleteMutation.isPending}
                          onClick={() => setDocumentToDelete(doc)}
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </div>
                    </Card>
                  );
                })}
              </div>
            )}
          </>
        )}
      </div>
      <ConfirmDialog
        open={Boolean(documentToDelete)}
        onOpenChange={(open) => {
          if (!open && !deleteMutation.isPending) setDocumentToDelete(null);
        }}
        variant="destructive"
        title="Удалить документ?"
        description={`Документ «${
          documentToDelete?.file_name ||
          DOCUMENT_TYPE_LABELS[documentToDelete?.document_type] ||
          'Файл'
        }» будет удалён из кабинета. Если файл был доступен на сервере, он также будет удалён из хранилища.`}
        confirmLabel={deleteMutation.isPending ? 'Удаление…' : 'Удалить'}
        disabled={deleteMutation.isPending}
        onConfirm={confirmDeleteDocument}
      />
    </div>
  );
}
