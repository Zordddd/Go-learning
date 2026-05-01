import React, { useState, useEffect } from 'react';
import {
    Box,
    List,
    ListItem,
    ListItemText,
    Checkbox,
    IconButton,
    TextField,
    Button,
    Paper,
    Typography,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Alert,
    CircularProgress,
    Grow,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import AddIcon from '@mui/icons-material/Add';
import CheckCircleOutlinedIcon from '@mui/icons-material/CheckCircleOutlined';
import RadioButtonUncheckedIcon from '@mui/icons-material/RadioButtonUnchecked';
import { motion, AnimatePresence } from 'framer-motion';
import { TaskService } from '../services/api';
import { Task, CreateTaskRequest, UpdateTaskRequest } from '../types/task';

export const TaskList: React.FC = () => {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [openDialog, setOpenDialog] = useState(false);
    const [editingTask, setEditingTask] = useState<Task | null>(null);
    const [formData, setFormData] = useState({ title: '', text: '' });

    useEffect(() => {
        loadTasks();
    }, []);

    const loadTasks = async () => {
        try {
            setLoading(true);
            const data = await TaskService.getTasks();
            setTasks(data);
            setError(null);
        } catch (err) {
            setError('Failed to load tasks');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = async () => {
        if (!formData.title.trim() || !formData.text.trim()) {
            setError('Title and text are required');
            return;
        }

        try {
            const newTask: CreateTaskRequest = {
                title: formData.title,
                text: formData.text,
            };
            await TaskService.createTask(newTask);
            await loadTasks();
            handleCloseDialog();
        } catch (err) {
            setError('Failed to create task');
            console.error(err);
        }
    };

    const handleUpdate = async () => {
        if (!editingTask) return;

        try {
            const updateData: UpdateTaskRequest = {
                id: editingTask.id,
                title: formData.title,
                text: formData.text,
                completed: editingTask.completed,
            };
            await TaskService.updateTask(updateData);
            await loadTasks();
            handleCloseDialog();
        } catch (err) {
            setError('Failed to update task');
            console.error(err);
        }
    };

    const handleToggleComplete = async (task: Task) => {
        try {
            const updateData: UpdateTaskRequest = {
                id: task.id,
                title: task.title,
                text: task.text,
                completed: !task.completed,
            };
            await TaskService.updateTask(updateData);
            await loadTasks();
        } catch (err) {
            setError('Failed to update task');
            console.error(err);
        }
    };

    const handleDelete = async (id: number) => {
        if (window.confirm('Are you sure you want to delete this task?')) {
            try {
                await TaskService.deleteTask(id);
                await loadTasks();
            } catch (err) {
                setError('Failed to delete task');
                console.error(err);
            }
        }
    };

    const handleEditClick = (task: Task) => {
        setEditingTask(task);
        setFormData({ title: task.title, text: task.text });
        setOpenDialog(true);
    };

    const handleCloseDialog = () => {
        setOpenDialog(false);
        setEditingTask(null);
        setFormData({ title: '', text: '' });
    };

    if (loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
                <CircularProgress size={60} thickness={4} sx={{ color: '#3b82f6' }} />
            </Box>
        );
    }

    return (
        <Box sx={{ maxWidth: 900, margin: 'auto', padding: 4 }}>
            <motion.div
                initial={{ opacity: 0, y: -20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6 }}
            >
                <Paper
                    elevation={0}
                    sx={{
                        background: 'rgba(15, 23, 42, 0.7)',
                        backdropFilter: 'blur(10px)',
                        borderRadius: '32px',
                        padding: 4,
                        border: '1px solid rgba(59, 130, 246, 0.2)',
                        boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
                    }}
                >
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                        <Typography
                            variant="h3"
                            component="h1"
                            sx={{
                                fontWeight: 700,
                                background: 'linear-gradient(135deg, #60a5fa, #a78bfa)',
                                backgroundClip: 'text',
                                WebkitBackgroundClip: 'text',
                                color: 'transparent',
                                textShadow: '0 0 20px rgba(96, 165, 250, 0.3)',
                            }}
                        >
                            Task Manager
                        </Typography>
                        <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
                            <Button
                                variant="contained"
                                startIcon={<AddIcon />}
                                onClick={() => setOpenDialog(true)}
                                sx={{
                                    background: 'linear-gradient(135deg, #3b82f6, #2563eb)',
                                    borderRadius: '40px',
                                    padding: '10px 24px',
                                    fontSize: '1rem',
                                    fontWeight: 600,
                                    textTransform: 'none',
                                    boxShadow: '0 4px 14px rgba(59, 130, 246, 0.4)',
                                    transition: 'all 0.3s ease',
                                    '&:hover': {
                                        background: 'linear-gradient(135deg, #2563eb, #1d4ed8)',
                                        transform: 'translateY(-2px)',
                                        boxShadow: '0 6px 20px rgba(59, 130, 246, 0.5)',
                                    },
                                }}
                            >
                                Add Task
                            </Button>
                        </motion.div>
                    </Box>

                    {error && (
                        <Alert
                            severity="error"
                            onClose={() => setError(null)}
                            sx={{
                                mb: 3,
                                borderRadius: '16px',
                                background: 'rgba(220, 38, 38, 0.1)',
                                backdropFilter: 'blur(8px)',
                                border: '1px solid rgba(220, 38, 38, 0.3)',
                            }}
                        >
                            {error}
                        </Alert>
                    )}

                    <List sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                        <AnimatePresence>
                            {tasks.length === 0 ? (
                                <motion.div
                                    initial={{ opacity: 0 }}
                                    animate={{ opacity: 1 }}
                                    exit={{ opacity: 0 }}
                                >
                                    <Typography
                                        align="center"
                                        sx={{
                                            py: 8,
                                            color: '#94a3b8',
                                            fontSize: '1.1rem',
                                        }}
                                    >
                                        ✨ No tasks yet. Create your first task!
                                    </Typography>
                                </motion.div>
                            ) : (
                                tasks.map((task, index) => (
                                    <Grow in timeout={300 + index * 100} key={task.id}>
                                        <motion.div
                                            initial={{ opacity: 0, x: -20 }}
                                            animate={{ opacity: 1, x: 0 }}
                                            exit={{ opacity: 0, x: 20 }}
                                            transition={{ duration: 0.3 }}
                                            whileHover={{ scale: 1.01 }}
                                        >
                                            <Paper
                                                elevation={0}
                                                sx={{
                                                    background: task.completed
                                                        ? 'rgba(34, 197, 94, 0.05)'
                                                        : 'rgba(30, 41, 59, 0.5)',
                                                    backdropFilter: 'blur(8px)',
                                                    borderRadius: '20px',
                                                    border: '1px solid rgba(59, 130, 246, 0.15)',
                                                    transition: 'all 0.3s ease',
                                                    '&:hover': {
                                                        borderColor: 'rgba(59, 130, 246, 0.5)',
                                                        boxShadow: '0 8px 32px rgba(0, 0, 0, 0.2)',
                                                        background: task.completed
                                                            ? 'rgba(34, 197, 94, 0.08)'
                                                            : 'rgba(30, 41, 59, 0.7)',
                                                    },
                                                }}
                                            >
                                                <ListItem sx={{ py: 2, px: 3 }}>
                                                    <motion.div whileTap={{ scale: 0.9 }}>
                                                        <Checkbox
                                                            checked={task.completed}
                                                            onChange={() => handleToggleComplete(task)}
                                                            icon={<RadioButtonUncheckedIcon sx={{ color: '#60a5fa' }} />}
                                                            checkedIcon={<CheckCircleOutlinedIcon sx={{ color: '#22c55e' }} />}
                                                            sx={{
                                                                '&.Mui-checked': {
                                                                    color: '#22c55e',
                                                                },
                                                            }}
                                                        />
                                                    </motion.div>
                                                    <ListItemText
                                                        primary={
                                                            <Typography
                                                                variant="h6"
                                                                sx={{
                                                                    textDecoration: task.completed ? 'line-through' : 'none',
                                                                    color: task.completed ? '#94a3b8' : '#f1f5f9',
                                                                    fontWeight: 500,
                                                                    transition: 'all 0.2s',
                                                                }}
                                                            >
                                                                {task.title}
                                                            </Typography>
                                                        }
                                                        secondary={
                                                            <Box sx={{ mt: 1 }}>
                                                                <Typography variant="body2" sx={{ color: '#cbd5e1', mb: 0.5 }}>
                                                                    {task.text}
                                                                </Typography>
                                                                <Typography variant="caption" sx={{ color: '#64748b' }}>
                                                                    {new Date(task.created_at).toLocaleString()}
                                                                </Typography>
                                                            </Box>
                                                        }
                                                    />
                                                    <Box sx={{ display: 'flex', gap: 1 }}>
                                                        <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
                                                            <IconButton
                                                                onClick={() => handleEditClick(task)}
                                                                sx={{
                                                                    color: '#60a5fa',
                                                                    '&:hover': {
                                                                        background: 'rgba(96, 165, 250, 0.1)',
                                                                    },
                                                                }}
                                                            >
                                                                <EditIcon />
                                                            </IconButton>
                                                        </motion.div>
                                                        <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
                                                            <IconButton
                                                                onClick={() => handleDelete(task.id)}
                                                                sx={{
                                                                    color: '#f87171',
                                                                    '&:hover': {
                                                                        background: 'rgba(248, 113, 113, 0.1)',
                                                                    },
                                                                }}
                                                            >
                                                                <DeleteIcon />
                                                            </IconButton>
                                                        </motion.div>
                                                    </Box>
                                                </ListItem>
                                            </Paper>
                                        </motion.div>
                                    </Grow>
                                ))
                            )}
                        </AnimatePresence>
                    </List>
                </Paper>
            </motion.div>

            <Dialog
                open={openDialog}
                onClose={handleCloseDialog}
                maxWidth="sm"
                fullWidth
                slotProps={{
                    paper: {
                        sx: {
                            background: 'rgba(15, 23, 42, 0.95)',
                            backdropFilter: 'blur(20px)',
                            borderRadius: '32px',
                            border: '1px solid rgba(59, 130, 246, 0.3)',
                            boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.5)',
                        },
                    },
                }}
            >
                <DialogTitle
                    sx={{
                        fontSize: '1.8rem',
                        fontWeight: 600,
                        background: 'linear-gradient(135deg, #60a5fa, #a78bfa)',
                        backgroundClip: 'text',
                        WebkitBackgroundClip: 'text',
                        color: 'transparent',
                        pb: 1,
                    }}
                >
                    {editingTask ? 'Edit Task' : 'Create New Task'}
                </DialogTitle>
                <DialogContent>
                    <TextField
                        autoFocus
                        margin="dense"
                        label="Title"
                        fullWidth
                        variant="outlined"
                        value={formData.title}
                        onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                        sx={{
                            mb: 3,
                            mt: 2,
                            '& .MuiOutlinedInput-root': {
                                borderRadius: '16px',
                                background: 'rgba(30, 41, 59, 0.6)',
                                '& fieldset': {
                                    borderColor: 'rgba(59, 130, 246, 0.3)',
                                },
                                '&:hover fieldset': {
                                    borderColor: '#3b82f6',
                                },
                            },
                            '& .MuiInputLabel-root': {
                                color: '#94a3b8',
                            },
                            '& .MuiInputBase-input': {
                                color: '#f1f5f9',
                            },
                        }}
                    />
                    <TextField
                        margin="dense"
                        label="Description"
                        fullWidth
                        multiline
                        rows={4}
                        variant="outlined"
                        value={formData.text}
                        onChange={(e) => setFormData({ ...formData, text: e.target.value })}
                        sx={{
                            '& .MuiOutlinedInput-root': {
                                borderRadius: '16px',
                                background: 'rgba(30, 41, 59, 0.6)',
                                '& fieldset': {
                                    borderColor: 'rgba(59, 130, 246, 0.3)',
                                },
                                '&:hover fieldset': {
                                    borderColor: '#3b82f6',
                                },
                            },
                            '& .MuiInputLabel-root': {
                                color: '#94a3b8',
                            },
                            '& .MuiInputBase-input': {
                                color: '#f1f5f9',
                            },
                        }}
                    />
                </DialogContent>
                <DialogActions sx={{ p: 3, gap: 2 }}>
                    <Button
                        onClick={handleCloseDialog}
                        sx={{
                            borderRadius: '40px',
                            padding: '8px 24px',
                            textTransform: 'none',
                            color: '#94a3b8',
                            '&:hover': {
                                background: 'rgba(255, 255, 255, 0.05)',
                            },
                        }}
                    >
                        Cancel
                    </Button>
                    <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                        <Button
                            onClick={editingTask ? handleUpdate : handleCreate}
                            variant="contained"
                            sx={{
                                background: 'linear-gradient(135deg, #3b82f6, #2563eb)',
                                borderRadius: '40px',
                                padding: '8px 32px',
                                fontWeight: 600,
                                textTransform: 'none',
                                boxShadow: '0 4px 14px rgba(59, 130, 246, 0.4)',
                                '&:hover': {
                                    background: 'linear-gradient(135deg, #2563eb, #1d4ed8)',
                                    transform: 'translateY(-1px)',
                                },
                            }}
                        >
                            {editingTask ? 'Update' : 'Create'}
                        </Button>
                    </motion.div>
                </DialogActions>
            </Dialog>
        </Box>
    );
};