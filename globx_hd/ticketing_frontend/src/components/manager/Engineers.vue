<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-5xl mx-auto mt-8">
    <h1 class="text-2xl font-bold text-blue-800 mb-6">Engineers</h1>
    <DataTable :columns="columns" :rows="engineers" rowKey="id" :rowClickable="true" @row-click="openEngineerDetail" />
    <Modal :open="showDetail" title="Assigned Tickets" @close="closeDetail">
      <div v-if="selectedEngineer">
        <div class="font-semibold text-blue-700 mb-2">
          {{ selectedEngineer.first_name }} {{ selectedEngineer.last_name }} <span class="text-gray-500">({{ selectedEngineer.phone || 'No phone' }})</span>
        </div>
        <div v-if="tickets.length">
          <div v-for="(ticket, idx) in tickets" :key="ticket.id" class="border-b border-blue-100 py-2">
            <div class="flex items-center justify-between cursor-pointer hover:bg-blue-50 px-2 py-1 rounded" @click="toggleTicket(idx)">
              <span class="font-mono text-blue-900">{{ ticket.ticket_id }}</span>
              <svg :class="[expanded[idx] ? 'rotate-180' : '', 'transition-transform']" width="20" height="20" fill="none" viewBox="0 0 24 24"><path stroke="#2563eb" stroke-width="2" d="M6 9l6 6 6-6"/></svg>
            </div>
            <div v-if="expanded[idx]" class="bg-blue-50 rounded p-3 mt-1">
              <div><span class="font-semibold text-gray-700">Product:</span> {{ ticket.product?.product_name || '-' }}</div>
              <div><span class="font-semibold text-gray-700">Issue:</span> {{ ticket.issue?.issue_name || ticket.issue_id }}</div>
              <div><span class="font-semibold text-gray-700">Details:</span> {{ ticket.ticket_details }}</div>
              <div><span class="font-semibold text-gray-700">Status:</span> {{ ticket.ticket_status }}</div>
              <div><span class="font-semibold text-gray-700">Contact:</span> {{ ticket.contact?.first_name }} {{ ticket.contact?.last_name }}</div>
              <div><span class="font-semibold text-gray-700">Account:</span> {{ ticket.account?.account_name }}</div>
              <div><span class="font-semibold text-gray-700">Created:</span> {{ formatDate(ticket.created_at) }}</div>
            </div>
          </div>
        </div>
        <div v-else class="text-blue-700 mt-4">No assigned tickets.</div>
      </div>
    </Modal>
  </div>
</template>
<script setup>
import { ref } from 'vue';
import { fetchEngineersWithTickets, fetchEngineerAssignedTickets } from '@/api/engineers.js';
import { formatDateIST } from '@/utils/date';
import DataTable from '@/components/ui/DataTable.vue';
import Modal from '@/components/ui/Modal.vue';

const formatDate = formatDateIST;
const engineers = ref([]);
const columns = [
  { key: 'first_name', label: 'Name' },
  { key: 'phone', label: 'Phone' },
  { key: 'assigned_tickets_count', label: 'Assigned Tickets' }
];
const showDetail = ref(false);
const selectedEngineer = ref(null);
const tickets = ref([]);
const expanded = ref([]);

async function loadEngineers() {
  const res = await fetchEngineersWithTickets();
  console.log('Engineers API response:', res);
  engineers.value = res;
}
async function openEngineerDetail(engineer) {
  selectedEngineer.value = engineer;
  tickets.value = await fetchEngineerAssignedTickets(engineer.id);
  expanded.value = Array(tickets.value.length).fill(false);
  showDetail.value = true;
}
function closeDetail() {
  showDetail.value = false;
  selectedEngineer.value = null;
  tickets.value = [];
  expanded.value = [];
}
function toggleTicket(idx) {
  expanded.value[idx] = !expanded.value[idx];
}
loadEngineers();
</script>
<style scoped>
.rotate-180 { transform: rotate(180deg); }
</style>
