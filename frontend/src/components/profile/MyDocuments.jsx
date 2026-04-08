import React, { useState } from 'react';
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
import { FileText, Upload, Trash2, Download } from 'lucide-react';
import { getApiErrorMessage } from '@/lib/apiErrors';
import {
  documentService,
  DOCUMENT_TYPE_LABELS,
} from '@/services/documentService';
import { ErrorNotice } from '@/components/common/ErrorNotice';

const DOC_TYPES = Object.keys(DOCUMENT_TYPE_LABELS);

export default function MyDocuments({ orders = [], appointments = [] }) {
  const queryClient = useQueryClient();
  const [file, setFile] = useState(null);
  const [documentType, setDocumentType] = useState('order_contract');
  const [linkKind, setLinkKind] = useState('order');
  const [orderId, setOrderId] = useState('');
  const [appointmentId, setAppointmentId] = useState('');
  const [formError, setFormError] = useState(null);

  const { data: documents, isLoading } = useQuery({
    queryKey: ['my-documents'],
    queryFn: () => documentService.list(),
  });

  const uploadMutation = useMutation({
    mutationFn: (payload) => documentService.upload(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-documents'] });
      setFile(null);
      setFormError(null);
    },
    onError: (e) => {
      setFormError(getApiErrorMessage(e, 'Не удалось загрузить файл'));
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id) => documentService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-documents'] });
      setFormError(null);
    },
    onError: (e) => {
      setFormError(getApiErrorMessage(e, 'Не удалось удалить документ'));
    },
  });

  const handleUpload = (e) => {
    e.preventDefault();
    if (!file) {
      setFormError('Выберите файл');
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

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2].map((i) => (
          <Card key={i} className="p-6 animate-pulse h-24" />
        ))}
      </div>
    );
  }

  const list = Array.isArray(documents) ? documents : [];

  return (
    <div className="space-y-8">
      <Card className="rounded-2xl border-slate-200 p-6 shadow-sm">
        <h3 className="text-lg font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Upload className="w-5 h-5" />
          Загрузить документ
        </h3>
        <p className="text-sm text-slate-600 mb-4">
          Прикрепите документ к заказу или записи на ТО. Файл будет доступен в вашем личном кабинете.
        </p>
        <form onSubmit={handleUpload} className="space-y-4 max-w-xl">
          <ErrorNotice kind="form" message={formError} />
          <div className="space-y-2">
            <Label>Файл</Label>
            <input
              type="file"
              className="block w-full text-sm text-slate-600"
              onChange={(ev) => {
                setFile(ev.target.files?.[0] || null);
                setFormError(null);
              }}
            />
          </div>
          <div className="space-y-2">
            <Label>Тип документа</Label>
            <Select value={documentType} onValueChange={(v) => {
              setDocumentType(v);
              setFormError(null);
            }}>
              <SelectTrigger>
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
            <Select value={linkKind} onValueChange={(v) => {
              setLinkKind(v);
              setFormError(null);
            }}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="order">К заказу</SelectItem>
                <SelectItem value="appointment">К записи на ТО</SelectItem>
              </SelectContent>
            </Select>
          </div>
          {linkKind === 'order' ? (
            <div className="space-y-2">
              <Label>Заказ</Label>
              <Select value={orderId} onValueChange={(v) => {
                setOrderId(v);
                setFormError(null);
              }}>
                <SelectTrigger>
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
          ) : (
            <div className="space-y-2">
              <Label>Запись на ТО</Label>
              <Select value={appointmentId} onValueChange={(v) => {
                setAppointmentId(v);
                setFormError(null);
              }}>
                <SelectTrigger>
                  <SelectValue placeholder="Выберите запись" />
                </SelectTrigger>
                <SelectContent>
                  {appointments.map((a) => {
                    const id =
                      a.service_appointment_id?.toString?.() ||
                      a.service_appointment_id;
                    return (
                      <SelectItem key={id} value={id}>
                        Запись #{String(id).slice(0, 8)}…
                      </SelectItem>
                    );
                  })}
                </SelectContent>
              </Select>
            </div>
          )}
          <Button type="submit" disabled={uploadMutation.isPending}>
            {uploadMutation.isPending ? 'Загрузка…' : 'Загрузить'}
          </Button>
        </form>
      </Card>

      <div>
        <h3 className="text-lg font-semibold text-slate-900 mb-4">
          Мои документы
        </h3>
        {list.length === 0 ? (
          <EmptyState
            icon={FileText}
            title="Нет документов"
            description="Загрузите файл, привязанный к заказу или записи на сервис"
          />
        ) : (
          <div className="space-y-3">
            {list.map((doc) => {
              const id = doc.document_id?.toString?.() || doc.document_id;
              const name =
                doc.file_name || DOCUMENT_TYPE_LABELS[doc.document_type] || 'Файл';
              return (
                <Card
                  key={id}
                  className="rounded-xl border-slate-200 p-4 shadow-sm flex flex-wrap items-center justify-between gap-3"
                >
                  <div className="flex items-start gap-3 min-w-0">
                    <div className="w-10 h-10 rounded-lg bg-slate-100 flex items-center justify-center flex-shrink-0">
                      <FileText className="w-5 h-5 text-slate-600" />
                    </div>
                    <div className="min-w-0">
                      <p className="font-medium text-slate-900 truncate">{name}</p>
                      <p className="text-sm text-slate-500">
                        {DOCUMENT_TYPE_LABELS[doc.document_type] || doc.document_type}
                      </p>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        documentService.download(id, name).catch((err) => {
                          setFormError(getApiErrorMessage(err, 'Не удалось скачать файл'));
                        })
                      }
                    >
                      <Download className="w-4 h-4 mr-1" />
                      Скачать
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      className="text-red-600 hover:text-red-700"
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
      </div>
    </div>
  );
}
