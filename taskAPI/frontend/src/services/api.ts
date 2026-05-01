import axios from 'axios';
import { Task, CreateTaskRequest, UpdateTaskRequest } from '../types/task';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Добавляем токен авторизации если нужен
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const TaskService = {
    // Получить все задачи
    async getTasks(): Promise<Task[]> {
        const response = await api.get('/tasks');
        return response.data;
    },

    // Создать задачу
    async createTask(task: CreateTaskRequest): Promise<Task> {
        const response = await api.post('/tasks', task);
        return response.data;
    },

    // Обновить задачу
    async updateTask(task: UpdateTaskRequest): Promise<Task> {
        const response = await api.put('/tasks', task);
        return response.data.data;
    },

    // Удалить задачу
    async deleteTask(id: number): Promise<void> {
        await api.delete(`/tasks?id=${id}`);
    },
};