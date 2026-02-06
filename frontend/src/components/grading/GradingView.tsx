'use client';

import { useGrading } from '@/hooks/useGrading'
import { useSubmission } from '@/hooks/useSubmission'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { AnswerDisplay } from './AnswerDisplay'
import { AIReasoningPanel } from './AIReasoningPanel'
import { OverrideForm } from './OverrideForm'

export function GradingView({ submissionId }: { submissionId: string }) {
    const { grades, loading: gradesLoading, error: gradesError, applyOverride } = useGrading(submissionId)
    const { submission, loading: subLoading, error: subError } = useSubmission(submissionId)
    
    if (gradesLoading || subLoading) return <LoadingSpinner />
    if (gradesError || subError) return <div className="p-4 text-red-500">Error: {gradesError || subError}</div>

    return (
        <div className="space-y-12">
            {grades.map((grade) => {
                const answer = submission?.answers?.find(a => a.question_id === grade.question_id);
                return (
                    <div key={grade.id} className="border-b pb-8 last:border-0">
                        <h2 className="text-xl font-semibold mb-4">Question ID: {grade.question_id}</h2>
                        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                            <div className="space-y-4">
                                {answer ? (
                                    <AnswerDisplay answer={{
                                        text: answer.text,
                                        diagrams: answer.diagrams
                                    }} />
                                ) : (
                                    <div className="p-4 bg-gray-100 rounded">Answer content not found for this grade.</div>
                                )}
                            </div>
                            <div className="space-y-4">
                                <AIReasoningPanel 
                                    reasoning={grade.reasoning}
                                    confidence={grade.confidence}
                                />
                                <div className="p-4 bg-blue-50 rounded-lg">
                                    <p className="font-bold mb-2">AI Score: {grade.final_score} / {grade.max_score}</p>
                                    <p className="text-sm text-gray-600">Status: {grade.status}</p>
                                </div>
                                <OverrideForm 
                                    currentScore={grade.final_score}
                                    onSubmit={(score, reason) => applyOverride(grade.question_id, score, reason)}
                                />
                            </div>
                        </div>
                    </div>
                );
            })}
            {grades.length === 0 && (
                <div className="text-center p-12 bg-gray-50 rounded-xl">
                    <p className="text-gray-500">No grades found for this submission yet.</p>
                </div>
            )}
        </div>
    )
}
