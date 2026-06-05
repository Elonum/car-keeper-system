import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
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

const DOC_TYPES = Object.keys(DOCUMENT_TYPE_LABELS);

function docId(doc) {
  return doc.document_id?.toString?.() || doc.document_id;
}

export default function MyDocuments({ orders = [], appointments = [] }) {
  const queryClient = useQueryClient();
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

  const canLinkOrder = orders.length > 0;
  const canLinkAppointment = appointments.length > 0;
  const canUpload = canLinkOrder || canLinkAppointment;

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

  const typeFilterOptions = useMemo(() => {
    const head = [{ value: 'all', label: 'Все типы' }];
    const rest = DOC_TYPES.map((t) => ({
      value: t,
      label: DOCUMENT_TYPE_LABELS[t],
    }));
    return [...head, ...rest];
  }, []);

  const filteredList = useMemo(() => {
    const q = search.trim().toLowerCase();
    return list.filter((doc) => {
      if (typeFilter !== 'all' && doc.document_type !== typeFilter) return false;
      if (!q) return true;
      const name = (doc.file_name || DOCUMENT_TYPE_LABELS[doc.document_type] || '').toLowerCase();
      const typeLabel = (DOCUMENT_TYPE_LABELS[doc.document_type] || doc.document_type || '').toLowerCase();
      const id = String(docId(doc) || '').toLowerCase();
      return name.includes(q) || typeLabel.includes(q) || id.includes(q);
    });
  }, [list, search, typeFilter]);

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
                          Заказ #{String(id).slice(0, 8)}…
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
                          Запись #{String(id).slice(0, 8)}…
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
        <h3 className="text-lg font-semibold text-slate-900">Мои документы</h3>
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
              searchPlaceholder="Поиск: имя файла, тип, номер…"
              statusValue={typeFilter}
              onStatusChange={setTypeFilter}
              statusOptions={typeFilterOptions}
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
                  const created =
                    doc.created_at &&
                    format(new Date(doc.created_at), 'd MMM yyyy, HH:mm', { locale: ru });

                  return (
                    <Card
                      key={id}
                      className="rounded-2xl border-slate-200 p-4 shadow-sm flex flex-wrap items-center justify-between gap-4"
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
                      <div className="flex gap-2 flex-shrink-0">
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
                          onClick={() => {
                            if (window.confirm('Удалить документ?')) {
                              deleteMutation.mutate(id);
                            }
                          }}
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
    </div>
  );
}
