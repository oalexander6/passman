import { ThemeProvider } from '@/lib/theme-provider';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import axios from 'axios';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';

import Navbar from '@/components/navbar';

import Dashboard from './components/dashboard';
import './index.css';

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: '/login',
    element: <p>Login page</p>
  },
  {
    path: '/register',
    element: <p>Register page</p>
  },
  {
    path: '/',
    element: <Navbar />,
    children: [
      {
        path: '',
        element: <Dashboard />
      }
    ]
  }
]);

axios.defaults.headers.common['X-XSRF-PROTECTION'] = 1;

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </ThemeProvider>
  </StrictMode>
);
