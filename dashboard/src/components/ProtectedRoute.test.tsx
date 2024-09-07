import { describe, it, expect, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import ProtectedRoute from './ProtectedRoute';
import { api } from '../api/apiClient';

// Mock the api module
vi.mock('../api/apiClient', () => ({
    api: {
        url: vi.fn().mockReturnThis(),
        get: vi.fn().mockReturnThis(),
        json: vi.fn(),
    },
}));

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
    },
});

const renderProtectedRoute = (isAuthenticated: boolean) => {
    (api.url('/api/user/profile').get().json as any).mockImplementation(() =>
        isAuthenticated ? Promise.resolve() : Promise.reject()
    );

    return render(
        <QueryClientProvider client={queryClient}>
            <MemoryRouter initialEntries={['/protected']}>
                <Routes>
                    <Route path="/protected" element={
                        <ProtectedRoute>
                            <div>Protected Content</div>
                        </ProtectedRoute>
                    } />
                    <Route path="/login" element={<div>Login Page</div>} />
                </Routes>
            </MemoryRouter>
        </QueryClientProvider>
    );
};

describe('ProtectedRoute', () => {
    it('shows loading state initially', async () => {
        renderProtectedRoute(true);
        expect(screen.getByText('Loading...')).toBeInTheDocument();
    });

    it('renders children when authenticated', async () => {
        renderProtectedRoute(true);
        await waitFor(() => {
            expect(screen.getByText('Protected Content')).toBeInTheDocument();
        });
    });

    it('redirects to login when not authenticated', async () => {
        renderProtectedRoute(false);
        await waitFor(() => {
            expect(screen.getByText('Login Page')).toBeInTheDocument();
        });
    });
});
