import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material';
import { TaskList } from './components/TaskList';
import './index.css'; // обязательно импортируйте глобальные стили

const theme = createTheme({
    palette: {
        mode: 'dark',
        primary: {
            main: '#3b82f6',
        },
        secondary: {
            main: '#a78bfa',
        },
        background: {
            default: 'transparent',
            paper: 'transparent',
        },
    },
    typography: {
        fontFamily: '"Inter", "Helvetica", "Arial", sans-serif',
    },
    components: {
        MuiCssBaseline: {
            styleOverrides: {
                body: {
                    backgroundColor: 'transparent',
                },
            },
        },
    },
});

function App() {
    return (
        <ThemeProvider theme={theme}>
            <Router>
                <Routes>
                    <Route path="/" element={<TaskList />} />
                </Routes>
            </Router>
        </ThemeProvider>
    );
}

export default App;