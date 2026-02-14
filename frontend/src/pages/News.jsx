import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { newsService } from '@/services/newsService';
import { Link } from 'react-router-dom';
import { createPageUrl } from '../utils';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import PageLoader from '../components/common/PageLoader';
import EmptyState from '../components/common/EmptyState';
import SectionHeader from '../components/common/SectionHeader';
import { Newspaper, Calendar, User, ArrowRight } from 'lucide-react';

export default function News() {
  const { data: news, isLoading } = useQuery({
    queryKey: ['news'],
    queryFn: () => newsService.getNews({ is_published: true }),
  });

  if (isLoading) return <PageLoader />;

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <SectionHeader 
          title="Новости"
          description="Последние новости и обновления автосалона"
        />

        {!news || news.length === 0 ? (
          <EmptyState 
            icon={Newspaper}
            title="Нет новостей"
            description="Скоро здесь появятся новости и обновления"
          />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {news.map(item => (
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
