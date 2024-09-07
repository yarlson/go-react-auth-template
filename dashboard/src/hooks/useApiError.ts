import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { isAuthError } from '../api/apiClient';

export const useApiError = (error: unknown) => {
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthError(error)) {
      navigate('/login');
    }
  }, [error, navigate]);

  if (error instanceof Error) {
    if (error.message.startsWith('Offline:')) {
      return 'You are currently offline. Some features may be unavailable.';
    }
    if (error.message.startsWith('Server Error:')) {
      return 'A server error occurred. Please try again later.';
    }
    return error.message;
  }

  return 'An unknown error occurred';
};
