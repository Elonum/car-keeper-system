import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { useSearchParams, Link } from 'react-router-dom';
import { newsService } from '@/services/newsService';
import { createPageUrl } from '../utils';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import PageLoader from '../components/common/PageLoader';
import { Button } from "@/components/ui/button";
import { ArrowLeft, Calendar, User } from 'lucide-react';
import ReactMarkdown from 'react-markdown';

export default function NewsDetail() {
  const [searchParams] = useSearchParams();
  const newsId = searchParams.get('id');

  const { data: news, isLoading } = useQuery({
    queryKey: ['news', newsId],
    queryFn: () => newsService.getNewsById(newsId),
    enabled: !!newsId,
  });

  if (isLoading) return <PageLoader />;
  if (!news) {
    return (
      <div className="min-h-screen bg-slate-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-slate-900 mb-2">Новость не найдена</h2>
          <Link to={createPageUrl("News")}>
            <Button variant="outline" className="mt-4 gap-2">
              <ArrowLeft className="w-4 h-4" />
              Вернуться к новостям
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Link to={createPageUrl("News")}>
          <Button variant="ghost" className="mb-6 gap-2 -ml-3">
            <ArrowLeft className="w-4 h-4" />
            Все новости
          </Button>
        </Link>

        <article className="bg-white rounded-2xl shadow-sm overflow-hidden">
          {news.image_url && (
            <div className="relative aspect-[21/9] overflow-hidden bg-slate-100">
              <img
                src={news.image_url}
                alt={news.title}
                className="w-full h-full object-cover"
              />
            </div>
          )}

          <div className="p-8 md:p-12">
            <h1 className="text-3xl md:text-4xl font-bold text-slate-900 mb-4 leading-tight">
              {news.title}
            </h1>

            <div className="flex flex-wrap items-center gap-4 pb-6 mb-6 border-b border-slate-100">
              {news.published_at && (
                <div className="flex items-center gap-2 text-sm text-slate-500">
                  <Calendar className="w-4 h-4" />
                  <span>
                    {format(new Date(news.published_at), 'd MMMM yyyy', { locale: ru })}
                  </span>
                </div>
              )}
              {news.author_name && (
                <div className="flex items-center gap-2 text-sm text-slate-500">
                  <User className="w-4 h-4" />
                  <span>{news.author_name}</span>
                </div>
              )}
            </div>

            <div className="prose prose-slate max-w-none prose-headings:font-bold prose-a:text-blue-600 prose-img:rounded-xl">
              <ReactMarkdown>{news.content}</ReactMarkdown>
            </div>
          </div>
        </article>
      </div>
    </div>
  );
}
