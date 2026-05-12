import api from './api'

export default {
  /**
   * Get audit logs with filters and pagination
   * @param {Object} params - Filter parameters
   * @returns {Promise} Audit logs with pagination metadata
   */
  async getAuditLogs(params = {}) {
    const response = await api.get('/manager/audit-logs', { params })
    return response.data
  },

  /**
   * Get a single audit log by ID
   * @param {number} id - Audit log ID
   * @returns {Promise} Audit log details
   */
  async getAuditLog(id) {
    const response = await api.get(`/manager/audit-logs/${id}`)
    return response.data
  },

  /**
   * Get audit log statistics
   * @returns {Promise} Statistics for last 30 days
   */
  async getStats() {
    const response = await api.get('/manager/audit-logs/stats')
    return response.data
  },

  /**
   * Get recent audit logs (last 30 days)
   * @param {number} limit - Maximum number of logs to retrieve
   * @returns {Promise} Recent audit logs
   */
  async getRecent(limit = 100) {
    const response = await api.get('/manager/audit-logs/recent', { params: { limit } })
    return response.data
  },

  /**
   * Get critical severity audit logs
   * @param {number} limit - Maximum number of logs to retrieve
   * @returns {Promise} Critical audit logs
   */
  async getCritical(limit = 100) {
    const response = await api.get('/manager/audit-logs/critical', { params: { limit } })
    return response.data
  },

  /**
   * Get failed operation audit logs
   * @param {number} limit - Maximum number of logs to retrieve
   * @returns {Promise} Failed audit logs
   */
  async getFailed(limit = 100) {
    const response = await api.get('/manager/audit-logs/failed', { params: { limit } })
    return response.data
  },

  /**
   * Get audit logs for a specific entity
   * @param {string} entityType - Entity type
   * @param {number} entityId - Entity ID
   * @param {number} limit - Maximum number of logs to retrieve
   * @returns {Promise} Entity audit logs
   */
  async getByEntity(entityType, entityId, limit = 100) {
    const response = await api.get(`/manager/audit-logs/entity/${entityType}/${entityId}`, { params: { limit } })
    return response.data
  },

  /**
   * Get audit logs for a specific actor
   * @param {string} actorType - Actor type
   * @param {number} actorId - Actor ID
   * @param {number} limit - Maximum number of logs to retrieve
   * @returns {Promise} Actor audit logs
   */
  async getByActor(actorType, actorId, limit = 100) {
    const response = await api.get(`/manager/audit-logs/actor/${actorType}/${actorId}`, { params: { limit } })
    return response.data
  }
}
