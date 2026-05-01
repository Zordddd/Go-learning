export interface Task {
    id: number;
    title: string;
    text: string;
    completed: boolean;
    created_at: string;
    updated_at: string;
}

export interface CreateTaskRequest {
    title: string;
    text: string;
}

export interface UpdateTaskRequest {
    id: number;
    title: string;
    text: string;
    completed: boolean;
}