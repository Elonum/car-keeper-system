import React, { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/lib/AuthContext';
import { profileService } from '@/services/profileService';
import { Button } from '@/components/ui/button';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import UserCarsList from './UserCarsList';
import AddUserCarDialog from './AddUserCarDialog';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { Plus } from 'lucide-react';
import { ErrorNotice } from '../common/ErrorNotice';

const GARAGE_QUERY_KEYS = [['my-cars'], ['userCars']];

export default function UserGarageSection() {
  const { isAuthenticated } = useAuth();
  const queryClient = useQueryClient();
  const [addOpen, setAddOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState(null);
  const [garageError, setGarageError] = useState(null);

  const { data: cars = [], isLoading } = useQuery({
    queryKey: ['my-cars'],
    queryFn: () => profileService.getUserCars(),
    enabled: isAuthenticated,
  });

  const deleteMutation = useMutation({
    mutationFn: (id) => profileService.deleteUserCar(id),
    onSuccess: () => {
      GARAGE_QUERY_KEYS.forEach((key) => queryClient.invalidateQueries({ queryKey: key }));
      setDeleteTarget(null);
      setGarageError(null);
    },
    onError: (e) => {
      setGarageError(getApiErrorMessage(e, 'Не удалось удалить автомобиль'));
    },
  });

  const list = Array.isArray(cars) ? cars : [];

  return (
    <div className="space-y-6">
      <ErrorNotice kind="server" message={garageError} />
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h3 className="text-lg font-semibold text-slate-900">Гараж</h3>
          <p className="mt-1 text-sm text-slate-600 max-w-xl">
            Автомобили привязаны к вашему аккаунту: по ним можно записаться на ТО. VIN уникален в системе.
            Удаление машины также удалит связанные с ней записи на обслуживание.
          </p>
        </div>
        <Button
          type="button"
          onClick={() => setAddOpen(true)}
          className="shrink-0 gap-2 bg-slate-900 hover:bg-slate-800"
        >
          <Plus className="h-4 w-4" />
          Добавить автомобиль
        </Button>
      </div>

      <UserCarsList
        cars={list}
        isLoading={isLoading}
        onRequestDelete={(car) => setDeleteTarget(car)}
        deletingId={
          deleteMutation.isPending && deleteMutation.variables
            ? deleteMutation.variables
            : null
        }
      />

      <AddUserCarDialog open={addOpen} onOpenChange={setAddOpen} />

      <AlertDialog open={!!deleteTarget} onOpenChange={(o) => !o && setDeleteTarget(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Удалить автомобиль?</AlertDialogTitle>
            <AlertDialogDescription>
              {deleteTarget && (
                <>
                  Будет удалена запись{' '}
                  <span className="font-medium text-slate-800">
                    {deleteTarget.brand_name} {deleteTarget.model_name}, VIN {deleteTarget.vin}
                  </span>
                  . Все записи на ТО по этому автомобилю будут удалены без восстановления.
                </>
              )}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <Button
              type="button"
              variant="destructive"
              disabled={deleteMutation.isPending}
              onClick={() => {
                if (!deleteTarget) return;
                const id = deleteTarget.user_car_id || deleteTarget.id;
                if (id) deleteMutation.mutate(id);
              }}
            >
              {deleteMutation.isPending ? 'Удаление…' : 'Удалить'}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
