import './App.css';
import axios from 'axios';
import Example from './components/Example';
import { useEffect, useState } from 'react';

function App() {
    axios.defaults.headers.common['X-XSRF-PROTECTION'] = 1;
    axios.defaults.baseURL = 'http://localhost:8000'

    const [msg, setMsg] = useState(null);

    useEffect(() => {
        axios.get('/api')
            .then((res => {
                setMsg(res.data);
            }))
            .catch();
    }, []);

    return (
        <>
            {msg}
            <Example></Example>
        </>
    );
}

export default App;
