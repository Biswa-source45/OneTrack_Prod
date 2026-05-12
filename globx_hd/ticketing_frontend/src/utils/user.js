/**
 * Utility functions for user data formatting
 */

/**
 * Format user's full name, handling cases where first_name already contains both names
 * @param {Object} user - User object with first_name and last_name
 * @returns {string} Properly formatted full name
 */
export const formatUserName = (user) => {
  if (!user) return 'Unknown User'
  
  const firstName = user.first_name || ''
  const lastName = user.last_name || ''
  
  // Check if first_name already contains both names (common data issue)
  if (firstName.includes(' ') || !lastName) {
    return firstName.trim()  // Use only first_name if it already contains full name
  }
  
  return `${firstName} ${lastName}`.trim()
}

/**
 * Get user initials for avatars
 * @param {Object} user - User object with first_name and last_name
 * @returns {string} User initials (e.g., "JD")
 */
export const getUserInitials = (user) => {
  if (!user) return '?'
  
  const firstName = user.first_name || ''
  const lastName = user.last_name || ''
  
  // If first_name contains space, use first letter of each word
  if (firstName.includes(' ')) {
    const names = firstName.trim().split(' ')
    return (names[0].charAt(0) + (names[1]?.charAt(0) || '')).toUpperCase()
  }
  
  return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase() || '?'
}
