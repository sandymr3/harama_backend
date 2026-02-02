import React from 'react'

export const AnswerDisplay = ({ answer }: { answer: any }) => (
    <div className="space-y-4">
        <div className="p-4 border rounded bg-gray-50">
            <h3 className="text-sm font-semibold text-gray-500 uppercase">Student Answer</h3>
            <p className="mt-2 whitespace-pre-wrap">{answer.text}</p>
        </div>
        
        {answer.diagrams && answer.diagrams.length > 0 && (
            <div className="grid grid-cols-2 gap-2">
                {answer.diagrams.map((url: string, idx: number) => (
                    <div key={idx} className="border rounded overflow-hidden">
                        <img src={url} alt={`Diagram ${idx + 1}`} className="w-full h-auto" />
                        <div className="bg-gray-100 p-1 text-xs text-center">Diagram {idx + 1}</div>
                    </div>
                ))}
            </div>
        )}
    </div>
)
