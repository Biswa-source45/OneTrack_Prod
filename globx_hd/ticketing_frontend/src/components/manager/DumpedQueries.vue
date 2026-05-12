<template>
  <div class="p-6 bg-white rounded-lg shadow-md max-w-7xl mx-auto mt-8">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-blue-800">Dumped Queries</h1>
      <div class="flex gap-2">
        <select v-model="statusFilter" @change="loadQueries" class="border rounded px-3 py-2 text-sm">
          <option value="">All Status</option>
          <option value="OPEN">Open</option>
          <option value="RESOLVED">Resolved</option>
          <option value="IGNORED">Ignored</option>
        </select>
        <button 
          @click="loadQueries" 
          class="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          Refresh
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-500">Loading...</div>

    <div v-else-if="queries.length === 0" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No dumped queries found</h3>
      <p class="mt-1 text-sm text-gray-500">Failed email parses will appear here.</p>
    </div>

    <div v-else class="space-y-4">
      <div v-for="query in queries" :key="query.id" 
           class="bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer"
           @click="openDetail(query)">
        
        <div class="p-4">
          <div class="flex justify-between items-start mb-2">
            <div>
              <h3 class="text-lg font-semibold text-gray-900 truncate pr-4">
                {{ query.subject || 'No Subject' }}
              </h3>
              <p class="text-sm text-gray-500">{{ query.sender_email }} <span v-if="query.sender_name">({{ query.sender_name }})</span></p>
            </div>
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border"
                  :class="getStatusClass(query.status)">
              {{ query.status }}
            </span>
          </div>
          
          <div class="mb-3">
             <p class="text-sm text-red-600 line-clamp-1">
               <span class="font-medium">Failure:</span> {{ query.failure_reason }}
             </p>
          </div>

          <div class="flex items-center justify-between text-xs text-gray-500 border-t pt-3 mt-3">
            <div class="flex gap-4">
               <span>ID: #{{ query.id }}</span>
               <span>N8N ID: {{ query.n8n_id }}</span>
            </div>
            <span>{{ formatDate(query.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-center mt-6 gap-2">
       <button 
         :disabled="currentPage === 1"
         @click="changePage(currentPage - 1)"
         class="px-3 py-1 rounded border disabled:opacity-50 hover:bg-gray-50"
       >
         Previous
       </button>
       <span class="px-3 py-1 text-gray-600">Page {{ currentPage }} of {{ totalPages }}</span>
       <button 
         :disabled="currentPage === totalPages"
         @click="changePage(currentPage + 1)"
         class="px-3 py-1 rounded border disabled:opacity-50 hover:bg-gray-50"
       >
         Next
       </button>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { fetchDumpedQueries } from '@/api/dumpedQueries';
import { formatDateIST } from '@/utils/date';

const router = useRouter();
const queries = ref([]);
const loading = ref(false);
const currentPage = ref(1);
const totalPages = ref(1);
const limit = 10;
const statusFilter = ref('OPEN');

const formatDate = formatDateIST;

async function loadQueries() {
  loading.value = true;
  try {
    const res = await fetchDumpedQueries(currentPage.value, limit, statusFilter.value);
    queries.value = res.data || [];
    const total = res.total || 0;
    totalPages.value = Math.ceil(total / limit);
  } catch (err) {
    console.error("Failed to load dumped queries", err);
  } finally {
    loading.value = false;
  }
}

function changePage(page) {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
    loadQueries();
  }
}

function getStatusClass(status) {
  switch (status) {
    case 'OPEN': return 'bg-blue-100 text-blue-800 border-blue-200';
    case 'RESOLVED': return 'bg-green-100 text-green-800 border-green-200';
    case 'IGNORED': return 'bg-gray-100 text-gray-800 border-gray-200';
    default: return 'bg-gray-100 text-gray-800';
  }
}

function openDetail(query) {
  router.push(`/manager/dumped-queries/${query.id}`);
}

onMounted(() => {
  loadQueries();
});
</script>
