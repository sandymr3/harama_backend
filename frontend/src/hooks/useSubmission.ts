'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { Submission } from '@/types';

export const useSubmission = (submissionId: string) => {
  const [submission, setSubmission] = useState<Submission | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSubmission = async () => {
      try {
        setLoading(true);
        // Assuming we'll add this to the API
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/submissions/${submissionId}`);
        if (!response.ok) throw new Error('Failed to fetch submission');
        const data = await response.json();
        setSubmission(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (submissionId) {
      fetchSubmission();
    }
  }, [submissionId]);

  return { submission, loading, error };
};
