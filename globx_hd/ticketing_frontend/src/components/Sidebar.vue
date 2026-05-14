<template>
  <aside :class="['bg-neutral-light min-h-full flex flex-col transition-all duration-300 overflow-hidden', sidebarOpen ? 'w-64 px-4 py-6 shadow-lg' : 'w-0 px-0 py-0 shadow-none']" aria-hidden="false">
    <nav class="flex flex-col h-full">
      <ul class="flex-1">
        <li v-for="item in navLinks" :key="item.path || item.label" class="mb-1">
          <router-link v-if="!item.children" :to="item.path" class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-teal-50 text-neutral-dark font-medium transition-all duration-200 backdrop-blur-sm" active-class="bg-brand-teal text-white border-l-4 border-teal-700">
            <span class="w-5 h-5 flex items-center justify-center">
              <component v-if="item.icon" :is="item.icon" class="w-5 h-5" :class="$route.path === item.path ? 'text-white' : 'text-brand-teal'" />
            </span>
            <span class="truncate">{{ item.label }}</span>
          </router-link>
          <div v-else class="w-full">
            <button @click="toggleMasterData" class="w-full flex items-center justify-between px-3 py-2 rounded-lg hover:bg-teal-50 text-neutral-dark font-medium transition-all duration-200">
              <span class="flex items-center gap-2">
                <span class="w-5 h-5 flex items-center justify-center">
                  <component v-if="item.icon" :is="item.icon" class="w-5 h-5 text-brand-teal" />
                </span>
                <span class="truncate">{{ item.label }}</span>
              </span>
              <svg :class="['transition-transform w-4 h-4 text-brand-teal', masterDataOpen ? 'rotate-180' : 'rotate-0']" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" /></svg>
            </button>
            <transition name="fade">
              <ul v-if="masterDataOpen" class="ml-7 mt-1">
                <li v-for="sub in item.children" :key="sub.path" class="mb-1">
                  <router-link :to="sub.path" class="block px-3 py-2 rounded-lg hover:bg-teal-50 text-neutral-dark font-medium transition-all duration-200" active-class="bg-teal-100 text-brand-teal">
                    {{ sub.label }}
                  </router-link>
                </li>
              </ul>
            </transition>
          </div>
        </li>
      </ul>
    </nav>
  </aside>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useAuthStore } from '../stores/auth';
import { BuildingOffice2Icon, UserIcon, Cog6ToothIcon } from '@heroicons/vue/24/outline';
import { PlusCircleIcon, TicketIcon, ClipboardDocumentListIcon, HomeIcon, UsersIcon, DocumentTextIcon } from '@heroicons/vue/24/outline';
const props = defineProps({ sidebarOpen: Boolean });
const masterDataOpen = ref(false);
function toggleMasterData() {
  masterDataOpen.value = !masterDataOpen.value;
}
watch(() => props.sidebarOpen, (val) => {
  if (!val) masterDataOpen.value = false;
});
const auth = useAuthStore();
const adminLinks = [
  { label: 'Accounts', path: '/accounts', icon: BuildingOffice2Icon },
  { label: 'Contacts', path: '/contacts', icon: UserIcon },
  { label: 'Users', path: '/users', icon: UserIcon },
  { label: 'Master Data', icon: Cog6ToothIcon, children: [
    { label: 'User Designations', path: '/master-data/user-designations' },
    { label: 'User Roles', path: '/master-data/user-roles' },
    { label: 'Contact Designations', path: '/master-data/contact-designations' },
    { label: 'Products', path: '/master-data/products' },
    { label: 'Issues', path: '/master-data/issues' },
  ] },
];
const managerLinks = [
  { label: 'Dashboard', path: '/dashboard', icon: HomeIcon },
  { label: 'Raise Ticket', path: '/manager/raise-ticket', icon: PlusCircleIcon },
  { label: 'Tickets', path: '/manager/tickets', icon: TicketIcon },
  { label: 'Dumped Queries', path: '/manager/dumped-queries', icon: ClipboardDocumentListIcon },
  { label: 'Create Task', path: '/manager/tasks/create', icon: PlusCircleIcon },
  { label: 'Tasks', path: '/manager/tasks', icon: ClipboardDocumentListIcon },
  { label: 'Engineers', path: '/manager/engineers', icon: UsersIcon },
  { label: 'Accounts', path: '/manager/accounts', icon: BuildingOffice2Icon },
  { label: 'Contacts', path: '/manager/contacts', icon: UserIcon },
  { label: 'Users', path: '/manager/users', icon: UserIcon },
  { label: 'Master Data', icon: Cog6ToothIcon, children: [
    { label: 'User Designations', path: '/manager/master-data/user-designations' },
    { label: 'User Roles', path: '/manager/master-data/user-roles' },
    { label: 'Contact Designations', path: '/manager/master-data/contact-designations' },
    { label: 'Products', path: '/manager/master-data/products' },
    { label: 'Issues', path: '/manager/master-data/issues' },
  ] },
  { label: 'Audit Logs', path: '/manager/audit-logs', icon: DocumentTextIcon },
];
const engineerLinks = [
  { label: 'Dashboard', path: '/dashboard', icon: HomeIcon },
  { label: 'Raise Ticket', path: '/engineer/raise-ticket', icon: PlusCircleIcon },
  { label: 'Assigned Tickets', path: '/engineer/tickets', icon: TicketIcon },
  { label: 'Assigned Tasks', path: '/engineer/tasks', icon: ClipboardDocumentListIcon },
];
const contactLinks = [
  { label: 'Raise Ticket', path: '/contacts/raise-ticket', icon: PlusCircleIcon },
  { label: 'My Tickets', path: '/contacts/my-tickets', icon: TicketIcon },
];
const navLinks = computed(() => {
  if (auth.userType === 'admin') return adminLinks;
  if (auth.userType === 'manager') return managerLinks;
  if (auth.userType === 'engineer') return engineerLinks;
  return contactLinks;
});
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
