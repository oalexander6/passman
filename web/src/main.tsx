import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import axios from 'axios';
import BaseLayout from './components/BaseLayout';

const queryClient = new QueryClient();

const router = createBrowserRouter([
    {
        path: '/login',
        element: <p>Login page</p>,
    },
    {
        path: '/register',
        element: <p>Register page</p>,
    },
    {
        path: '/',
        element: <BaseLayout />,
    },
]);

axios.defaults.headers.common['X-XSRF-PROTECTION'] = 1;

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <QueryClientProvider client={queryClient}>
            <RouterProvider router={router} />
        </QueryClientProvider>
    </StrictMode>
);
