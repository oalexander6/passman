import './App.css';
import axios from 'axios';
import Example from './components/Example';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();

function App() {
    axios.defaults.headers.common['X-XSRF-PROTECTION'] = 1;

    return (
        <QueryClientProvider client={queryClient}>
            <Example></Example>
        </QueryClientProvider>
    );
}

export default App;
