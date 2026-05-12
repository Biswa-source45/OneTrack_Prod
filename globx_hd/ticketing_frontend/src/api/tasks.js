import api from './api'

// Task CRUD operations
export function createTask(taskData) {
  return api.post('/manager/tasks', taskData).then(r => r.data)
}

export function fetchTasks() {
  return api.get('/manager/tasks').then(r => r.data)
}

export function fetchTask(id) {
  return api.get(`/manager/tasks/${id}`).then(r => r.data)
}

export function updateTask(id, taskData) {
  return api.put(`/manager/tasks/${id}`, taskData).then(r => r.data)
}

export function deleteTask(id) {
  return api.delete(`/manager/tasks/${id}`).then(r => r.data)
}

// Task comments
export const taskComments = {
  fetch: (taskId, params = {}) => {
    const query = new URLSearchParams(params).toString()
    return api.get(`/tasks/${taskId}/comments?${query}`).then(r => r.data)
  },
  
  create: (taskId, commentData) => {
    return api.post(`/tasks/${taskId}/comments`, commentData).then(r => r.data)
  },
  
  update: (taskId, commentId, commentData) => {
    return api.put(`/tasks/${taskId}/comments/${commentId}`, commentData).then(r => r.data)
  },
  
  delete: (taskId, commentId) => {
    return api.delete(`/tasks/${taskId}/comments/${commentId}`).then(r => r.data)
  }
}

// Task activities
export const taskActivities = {
  fetch: (taskId, params = {}) => {
    const query = new URLSearchParams(params).toString()
    return api.get(`/tasks/${taskId}/activities?${query}`).then(r => r.data)
  }
}

// Get task full details (task + comments + activities)
export async function fetchTaskFullDetails(id) {
  try {
    const [taskResponse, commentsResponse] = await Promise.all([
      fetchTask(id),
      taskComments.fetch(id, { limit: 50 })
    ])
    
    return {
      task: taskResponse.task,
      counts: {
        comments: commentsResponse.total || 0
      }
    }
  } catch (error) {
    console.error('Error fetching task details:', error)
    throw error
  }
}
