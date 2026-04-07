import React from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { newsService } from '@/services/newsService';
import { Link } from 'react-router-dom';
import { createPageUrl } from '../utils';
import { useAuth } from '@/lib/AuthContext';
import { PERMISSIONS, hasPermission } from '@/lib/authz';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { getApiErrorMessage } from '@/lib/apiErrors';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { Card } from "@/components/ui/card";
import PageLoader from '../components/common/PageLoader';
import EmptyState from '../components/common/EmptyState';
import SectionHeader from '../components/common/SectionHeader';
import { ErrorNotice } from '../components/common/ErrorNotice';
import { Newspaper, Calendar, User, ArrowRight } from 'lucide-react';

export default function News() {
  const qc = useQueryClient();
  const { user } = useAuth();
  const canManageNews = hasPermission(user?.role, PERMISSIONS.NEWS_MANAGE);
  const [scope, setScope] = React.useState('published');
  const [draftTitle, setDraftTitle] = React.useState('');
  const [draftContent, setDraftContent] = React.useState('');
  const [newsError, setNewsError] = React.useState(null);
  const { data: news, isLoading } = useQuery({
    queryKey: ['news', scope, canManageNews],
    queryFn: () =>
      newsService.getNews(
        canManageNews && scope !== 'published'
          ? { scope: scope === 'all' ? 'all' : 'unpublished' }
          : {}
      ),
  });
  const createMutation = useMutation({
    mutationFn: () => newsService.createNews({ title: draftTitle, content: draftContent }),
    onSuccess: () => {
      setDraftTitle('');
      setDraftContent('');
      setNewsError(null);
      qc.invalidateQueries({ queryKey: ['news'] });
    },
    onError: (e) => setNewsError(getApiErrorMessage(e, 'Не удалось создать новость')),
  });
  const publishMutation = useMutation({
    mutationFn: ({ id, publish } = {}) =>
      publish ? newsService.publishNews(id) : newsService.unpublishNews(id),
    onSuccess: () => {
      setNewsError(null);
      qc.invalidateQueries({ queryKey: ['news'] });
    },
    onError: (e) => setNewsError(getApiErrorMessage(e, 'Не удалось изменить публикацию')),
  });
  const deleteMutation = useMutation({
    mutationFn: (id = '') => newsService.deleteNews(id),
    onSuccess: () => {
      setNewsError(null);
      qc.invalidateQueries({ queryKey: ['news'] });
    },
    onError: (e) => setNewsError(getApiErrorMessage(e, 'Не удалось удалить новость')),
  });

  if (isLoading) return <PageLoader />;
  const list = Array.isArray(news) ? news : [];

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <SectionHeader 
          title="Новости"
          description="Последние новости и обновления автосалона"
        />
        <ErrorNotice kind="server" message={newsError} className="mb-6" />
        {canManageNews && (
          <Card className="p-5 mb-6 space-y-3">
            <div className="flex flex-wrap gap-2">
              <Button variant={scope === 'published' ? 'default' : 'outline'} size="sm" onClick={() => setScope('published')}>
                Опубликованные
              </Button>
              <Button variant={scope === 'unpublished' ? 'default' : 'outline'} size="sm" onClick={() => setScope('unpublished')}>
                Черновики
              </Button>
              <Button variant={scope === 'all' ? 'default' : 'outline'} size="sm" onClick={() => setScope('all')}>
                Все
              </Button>
            </div>
            <div className="grid md:grid-cols-2 gap-3">
              <Input placeholder="Заголовок новой новости" value={draftTitle} onChange={(e) => setDraftTitle(e.target.value)} />
              <Button
                onClick={() => createMutation.mutate()}
                disabled={createMutation.isPending || !draftTitle.trim() || !draftContent.trim()}
              >
                Создать черновик
              </Button>
            </div>
            <Textarea
              placeholder="Текст новости (markdown поддерживается)"
              value={draftContent}
              onChange={(e) => setDraftContent(e.target.value)}
              rows={5}
            />
          </Card>
        )}

        {list.length === 0 ? (
          <EmptyState 
            icon={Newspaper}
            title="Нет новостей"
            description="Скоро здесь появятся новости и обновления"
          />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {list.map(item => (
              <Link key={item.news_id || item.id} to={createPageUrl("NewsDetail") + `?id=${item.news_id || item.id}`}>
                <Card className="group overflow-hidden border-0 shadow-sm hover:shadow-xl transition-all duration-500 bg-white rounded-2xl h-full">
                  {/* Image */}
                  {item.image_url ? (
                    <div className="relative aspect-video overflow-hidden bg-slate-100">
                      <img
                        src={item.image_url}
                        alt={item.title}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700"
                      />
                    </div>
                  ) : (
                    <div className="relative aspect-video bg-gradient-to-br from-slate-100 to-slate-200 flex items-center justify-center">
                      <Newspaper className="w-12 h-12 text-slate-300" />
                    </div>
                  )}

                  {/* Content */}
                  <div className="p-5">
                    <h3 className="text-lg font-bold text-slate-900 mb-2 line-clamp-2 group-hover:text-blue-600 transition-colors">
                      {item.title}
                    </h3>
                    
                    {item.excerpt && (
                      <p className="text-sm text-slate-500 line-clamp-3 mb-4">
                        {item.excerpt}
                      </p>
                    )}

                    {/* Meta */}
                    <div className="flex flex-wrap items-center gap-3 text-xs text-slate-400">
                      {item.published_at && (
                        <div className="flex items-center gap-1">
                          <Calendar className="w-3.5 h-3.5" />
                          <span>
                            {format(new Date(item.published_at), 'd MMMM yyyy', { locale: ru })}
                          </span>
                        </div>
                      )}
                      {item.author_name && (
                        <div className="flex items-center gap-1">
                          <User className="w-3.5 h-3.5" />
                          <span>{item.author_name}</span>
                        </div>
                      )}
                    </div>

                    <div className="mt-4 flex items-center gap-1.5 text-sm font-medium text-blue-600 group-hover:gap-2.5 transition-all">
                      Читать далее
                      <ArrowRight className="w-4 h-4" />
                    </div>
                    {canManageNews && (
                      <div className="mt-3 flex gap-2">
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          onClick={(e) => {
                            e.preventDefault();
                            publishMutation.mutate({ id: item.news_id || item.id, publish: !item.is_published });
                          }}
                        >
                          {item.is_published ? 'Снять с публикации' : 'Опубликовать'}
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          className="text-red-600"
                          onClick={(e) => {
                            e.preventDefault();
                            deleteMutation.mutate(item.news_id || item.id);
                          }}
                        >
                          Удалить
                        </Button>
                      </div>
                    )}
                  </div>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
