<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header with Back Button -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-full px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center space-x-4">
            <button
              @click="goBack"
              class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              Back to List
            </button>
            <div class="h-6 w-px bg-gray-300"></div>
            <div v-if="query">
              <h1 class="text-xl font-semibold text-gray-900">Dumped Query #{{ query.id }}</h1>
              <p class="text-sm text-gray-500">{{ formatDate(query.created_at) }}</p>
            </div>
          </div>
          
          <div class="flex items-center space-x-3" v-if="query">
             <button 
              v-if="query.status !== 'IGNORED'"
              @click="updateStatus('IGNORED')"
              class="px-3 py-2 bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 text-sm font-medium shadow-sm"
            >
              Mark Ignored
            </button>
            <button 
              v-if="query.status !== 'RESOLVED'"
              @click="updateStatus('RESOLVED')"
              class="px-3 py-2 bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 text-sm font-medium shadow-sm"
            >
              Mark Resolved
            </button>
             <button 
              v-if="query.status !== 'OPEN'"
              @click="updateStatus('OPEN')"
              class="px-3 py-2 bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 text-sm font-medium shadow-sm"
            >
              Re-open
            </button>
             <button 
              @click="deleteQuery"
              class="px-3 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 text-sm font-medium shadow-sm"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <div v-else-if="error" class="p-8 text-center">
      <div class="text-red-600 mb-4">Error loading query</div>
      <p class="text-gray-600">{{ error }}</p>
    </div>

    <!-- Main Content Layout -->
    <div v-else-if="query" class="flex max-w-full">
      <!-- Left Sidebar - Properties Panel -->
      <div class="w-80 bg-white border-r border-gray-200 min-h-[calc(100vh-4rem)] p-4 space-y-6">
        
        <!-- Status -->
        <div>
          <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Status</h3>
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border"
                  :class="getStatusClass(query.status)">
              {{ query.status }}
          </span>
        </div>

        <!-- Failure Reason -->
        <div>
          <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Failure Reason</h3>
          <div class="text-sm text-red-700 bg-red-50 border border-red-100 p-3 rounded-md">
             {{ query.failure_reason }}
          </div>
        </div>

        <!-- Sender Info -->
        <div>
           <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Sender Information</h3>
           <div class="space-y-3 bg-gray-50 p-3 rounded-md border text-sm">
             <div>
               <span class="block text-gray-500 text-xs">Email</span>
               <span class="font-medium text-gray-900 break-all">{{ query.sender_email }}</span>
             </div>
             <div>
               <span class="block text-gray-500 text-xs">Name</span>
               <span class="font-medium text-gray-900">{{ query.sender_name || '-' }}</span>
             </div>
           </div>
        </div>
        
        <!-- Technical Info -->
        <div>
           <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Technical Info</h3>
           <div class="space-y-2 text-xs text-gray-500">
             <div class="flex justify-between">
               <span>ID:</span>
               <span class="font-mono">{{ query.id }}</span>
             </div>
             <div class="flex justify-between">
               <span>N8N ID:</span>
               <span class="font-mono">{{ query.n8n_id }}</span>
             </div>
           </div>
        </div>

      </div>

      <!-- Right Content Area -->
      <div class="flex-1 bg-white p-8">
        <div class="max-w-4xl mx-auto space-y-8">
          
          <!-- Subject -->
          <div>
            <h2 class="text-lg font-medium text-gray-900 mb-2">Subject</h2>
            <div class="bg-gray-50 border border-gray-200 rounded-md p-4 text-gray-900 font-medium">
              {{ query.subject || 'No Subject' }}
            </div>
          </div>

          <!-- Body -->
          <div>
            <h2 class="text-lg font-medium text-gray-900 mb-2">Email Content</h2>
            <div class="bg-gray-50 border border-gray-200 rounded-md p-4 text-gray-700 whitespace-pre-wrap font-mono text-sm leading-relaxed">
              {{ query.body }}
            </div>
          </div>

          <!-- AI Data -->
          <div>
            <h2 class="text-lg font-medium text-gray-900 mb-2">AI Extracted Data</h2>
            <div class="bg-gray-900 rounded-md p-4 overflow-x-auto">
               <pre class="text-xs text-green-400 font-mono">{{ formatJson(query.ai_extracted_data) }}</pre>
            </div>
          </div>
          
          <!-- Bottom Action -->
          <div class="pt-8 border-t border-gray-200 flex justify-center">
             <button 
               @click="convertToTicket"
               class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
             >
               <svg class="w-5 h-5 mr-2 -ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                 <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
               </svg>
               Convert to Ticket
             </button>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { fetchDumpedQuery, updateDumpedQueryStatus, deleteDumpedQuery } from '@/api/dumpedQueries';
import { formatDateIST } from '@/utils/date';

const route = useRoute();
const router = useRouter();
const queryId = route.params.id;
const query = ref(null);
const loading = ref(true);
const error = ref('');
const formatDate = formatDateIST;

async function loadQuery() {
  loading.value = true;
  error.value = '';
  try {
    const data = await fetchDumpedQuery(queryId);
    query.value = data;
  } catch (err) {
    console.error('Failed to load dumped query:', err);
    error.value = 'Failed to load details.';
  } finally {
    loading.value = false;
  }
}

function goBack() {
  router.push('/manager/dumped-queries');
}

function getStatusClass(status) {
  switch (status) {
    case 'OPEN': return 'bg-blue-100 text-blue-800 border-blue-200';
    case 'RESOLVED': return 'bg-green-100 text-green-800 border-green-200';
    case 'IGNORED': return 'bg-gray-100 text-gray-800 border-gray-200';
    default: return 'bg-gray-100 text-gray-800';
  }
}

function formatJson(jsonStr) {
  try {
    if (!jsonStr) return '{}';
    if (typeof jsonStr === 'object') return JSON.stringify(jsonStr, null, 2);
    const obj = JSON.parse(jsonStr);
    return JSON.stringify(obj, null, 2);
  } catch (e) {
    return jsonStr || '{}';
  }
}

async function updateStatus(newStatus) {
  if (!query.value) return;
  try {
    await updateDumpedQueryStatus(query.value.id, newStatus);
    query.value.status = newStatus;
  } catch (err) {
    alert("Failed to update status");
  }
}

async function deleteQuery() {
  if (!confirm("Are you sure you want to delete this query log?")) return;
  try {
    await deleteDumpedQuery(query.value.id);
    router.push('/manager/dumped-queries');
  } catch (err) {
    alert("Failed to delete query");
  }
}

function convertToTicket() {
  if (!query.value) return;
  
  router.push({
    path: '/manager/raise-ticket',
    query: {
      subject: query.value.subject || '',
      description: query.value.body || '',
      email: query.value.sender_email || '',
      name: query.value.sender_name || '',
      dump_id: query.value.id // Optional: Pass ID if we want to auto-resolve later
    }
  });
}

onMounted(() => {
  if (queryId) {
    loadQuery();
  }
});
</script>
