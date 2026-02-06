'use client';

import { useExams } from '@/hooks/useExams'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { ExamCard } from '@/components/exam/ExamCard'
import Link from 'next/link'

export default function ExamsPage() {
    const { exams, loading, error } = useExams()

    if (loading) return <LoadingSpinner />
    if (error) return <div className="p-4 text-red-500">Error: {error}</div>

    return (
        <div className="p-8">
            <div className="flex justify-between items-center mb-8">
                <h1 className="text-3xl font-bold">Exams</h1>
                <Link 
                    href="/exams/create" 
                    className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition"
                >
                    Create New Exam
                </Link>
            </div>
            
            {exams.length === 0 ? (
                <div className="text-center p-12 bg-gray-50 rounded-xl">
                    <p className="text-gray-500 mb-4">No exams found.</p>
                    <Link 
                        href="/exams/create" 
                        className="text-blue-600 hover:underline"
                    >
                        Create your first exam
                    </Link>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {exams.map((exam) => (
                        <ExamCard key={exam.id} exam={exam} />
                    ))}
                </div>
            )}
        </div>
    )
}
