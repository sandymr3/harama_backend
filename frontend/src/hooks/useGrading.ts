'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { FinalGrade } from '@/types';

export const useGrading = (submissionId: string) => {
  const [grades, setGrades] = useState<FinalGrade[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchGrades = async () => {
      try {
        setLoading(true);
        const data = await api.getGrades(submissionId);
        setGrades(data);
        setError(null);
      } catch (err: any) {
        setError(err.message || 'Failed to fetch grades');
      } finally {
        setLoading(false);
      }
    };

    if (submissionId) {
      fetchGrades();
    }
  }, [submissionId]);

  const applyOverride = async (questionId: string, score: number, reason: string) => {
    try {
      const updatedGrade = await api.captureOverride(submissionId, questionId, { score, reason });
      setGrades((prev) => 
        prev.map((g) => (g.question_id === questionId ? updatedGrade : g))
      );
      return updatedGrade;
    } catch (err: any) {
      throw new Error(err.message || 'Failed to apply override');
    }
  };

  return {
    grades,
    loading,
    error,
    applyOverride,
    refresh: async () => {
      const data = await api.getGrades(submissionId);
      setGrades(data);
    }
  };
};
