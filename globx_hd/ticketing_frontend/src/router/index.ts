import { createRouter, createWebHistory } from 'vue-router'


import LoginUser from '../components/LoginUser.vue';
import LoginContact from '../components/LoginContact.vue';
import ResetPasswordUser from '../components/ResetPasswordUser.vue';
import ResetPasswordContact from '../components/ResetPasswordContact.vue';
import Layout from '../components/Layout.vue';
import { useAuthStore } from '../stores/auth';
import { decodeJWT } from '../utils/jwt';

const publicPages = [
  '/login/user',
  '/login/contact',
  '/reset-password/user',
  '/reset-password/contact'
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/login/user', component: LoginUser },
    { path: '/login/contact', component: LoginContact },
    { path: '/reset-password/user', component: ResetPasswordUser },
    { path: '/reset-password/contact', component: ResetPasswordContact },
    {
      path: '/',
      component: Layout,
      children: [
        { path: 'dashboard', component: () => import('../components/Dashboard.vue'), meta: { requiresRole: 'admin' } },
        { path: 'accounts', component: () => import('../components/Accounts.vue'), meta: { requiresRole: 'admin' } },
        { path: 'contacts', component: () => import('../components/Contacts.vue'), meta: { requiresRole: 'admin' } },
        { path: 'users', component: () => import('../components/Users.vue'), meta: { requiresRole: 'admin' } },
        { path: 'master-data', component: () => import('../components/MasterData.vue'), meta: { requiresRole: 'admin' } },
        { path: 'master-data/user-designations', component: () => import('../components/masterdata/UserDesignations.vue'), meta: { requiresRole: 'admin' } },
        { path: 'master-data/user-roles', component: () => import('../components/masterdata/UserRoles.vue'), meta: { requiresRole: 'admin' } },
        { path: 'master-data/contact-designations', component: () => import('../components/masterdata/ContactDesignations.vue'), meta: { requiresRole: 'admin' } },
        { path: 'master-data/products', component: () => import('../components/masterdata/Products.vue'), meta: { requiresRole: 'admin' } },
        { path: 'master-data/products/new', component: () => import('../components/masterdata/ProductForm.vue'), meta: { requiresRole: 'admin' } },
        { path: 'master-data/products/:id/edit', component: () => import('../components/masterdata/ProductForm.vue'), meta: { requiresRole: 'admin' } },
        { path: 'master-data/issues', component: () => import('../components/masterdata/Issues.vue'), meta: { requiresRole: 'admin' } },
        // Contacts module child routes
        { path: 'contacts/raise-ticket', component: () => import('../components/contacts/RaiseTicket.vue'), meta: { requiresRole: 'contact' } },
        { path: 'contacts/my-tickets', component: () => import('../components/contacts/MyTickets.vue'), meta: { requiresRole: 'contact' } },
        { path: 'contacts/my-tickets/:id', component: () => import('../components/contacts/TicketDetail.vue'), meta: { requiresRole: 'contact' } },
                // Manager module
        { path: 'manager/dashboard', component: () => import('../components/manager/ManagerDashboard.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/raise-ticket', component: () => import('../components/manager/RaiseTicket.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/tickets', component: () => import('../components/manager/Tickets.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/tickets/:id', component: () => import('../components/shared/TicketDetailPage.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/tickets/:id/edit', component: () => import('../components/shared/EditTicketPage.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/tickets/:id/add-call', component: () => import('../components/shared/AddCallPage.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/engineers', component: () => import('../components/manager/Engineers.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/dumped-queries', component: () => import('../components/manager/DumpedQueries.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/dumped-queries/:id', component: () => import('../components/manager/DumpedQueryDetail.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/tasks', component: () => import('../components/manager/Tasks.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/tasks/create', component: () => import('../components/manager/CreateTask.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/tasks/:id', component: () => import('../components/shared/TaskDetailPage.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/tasks/:id/edit', component: () => import('../components/shared/EditTaskPage.vue'), meta: { requiresRole: 'manager' } },
        // Manager admin features
        { path: 'manager/accounts', component: () => import('../components/Accounts.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/contacts', component: () => import('../components/Contacts.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/users', component: () => import('../components/Users.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/master-data', component: () => import('../components/MasterData.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/master-data/user-designations', component: () => import('../components/masterdata/UserDesignations.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/master-data/user-roles', component: () => import('../components/masterdata/UserRoles.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/master-data/contact-designations', component: () => import('../components/masterdata/ContactDesignations.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/master-data/products', component: () => import('../components/masterdata/Products.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/master-data/products/new', component: () => import('../components/masterdata/ProductForm.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/master-data/products/:id/edit', component: () => import('../components/masterdata/ProductForm.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/master-data/issues', component: () => import('../components/masterdata/Issues.vue'), meta: { requiresRole: 'manager' } },
        { path: 'manager/audit-logs', component: () => import('../components/manager/AuditLogs.vue'), meta: { requiresRole: 'manager' } },
        // Engineer module
        { path: 'engineer/raise-ticket', component: () => import('../components/engineer/RaiseTicket.vue'), meta: { requiresRole: 'engineer' } },
        { path: 'engineer/tickets', component: () => import('../components/engineer/EngineerTickets.vue'), meta: { requiresRole: 'engineer' } },
        { path: 'engineer/tickets/:id', component: () => import('../components/engineer/EngineerTicketDetailPage.vue'), meta: { requiresRole: 'engineer' } },
        { path: 'engineer/tickets/:id/edit', component: () => import('../components/shared/EditTicketPage.vue'), meta: { requiresRole: 'engineer' } },
        { path: 'engineer/tickets/:id/add-call', component: () => import('../components/shared/AddCallPage.vue'), meta: { requiresRole: 'engineer' } },
        // Engineer task routes
        { path: 'engineer/tasks', component: () => import('../components/engineer/EngineerTasks.vue'), meta: { requiresRole: 'engineer' } },
        { path: 'engineer/tasks/:id', component: () => import('../components/engineer/EngineerTaskDetailPage.vue'), meta: { requiresRole: 'engineer' } },
        
        // Notifications (available to all authenticated users)
        { path: 'notifications', component: () => import('../components/shared/NotificationsPage.vue') },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/login/user' },
  ],
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore();
  if (!auth.token) auth.loadAuth();
  const token = auth.token;
  let isAuthenticated = false;
  let role = auth.userType;
  if (token)  {
    const claims = decodeJWT(token);
    if (claims && claims.exp && claims.exp * 1000 > Date.now()) {
      isAuthenticated = true;
      role = claims.role ? claims.role.toLowerCase().trim() : role;
    } else {
      auth.clearAuth();
      role = null;
    }
  }

  // If navigating to a public page, allow
  if (publicPages.includes(to.path)) {
    if (isAuthenticated && to.path.startsWith('/login')) {
      // Check if user needs to reset password first
      if (auth.firstLogin) {
        return next('/reset-password/user');
      }
      // Redirect to correct landing page based on role
      if (role === 'admin') return next('/dashboard');
      if (role === 'contact') return next('/contacts/my-tickets');
      if (role === 'manager') return next('/manager/dashboard');
      if (role === 'engineer') return next('/engineer/tickets');
    }
    return next();
  }

  // If not authenticated, redirect to login
  if (!isAuthenticated) {
    return next('/login/user');
  }

  // CRITICAL SECURITY CHECK: Force password reset for first-time users
  if (auth.firstLogin && !to.path.startsWith('/reset-password')) {
    return next('/reset-password/user');
  }

  // Role-based route protection
  const requiredRole = to.meta && to.meta.requiresRole;
  if (requiredRole && role !== requiredRole) {
    // Redirect to correct landing page
    if (role === 'admin') return next('/dashboard');
    if (role === 'contact') return next('/contacts/my-tickets');
    if (role === 'manager') return next('/manager/dashboard');
    if (role === 'engineer') return next('/engineer/tickets');
    return next('/login/user');
  }

  return next();
});

export default router;
