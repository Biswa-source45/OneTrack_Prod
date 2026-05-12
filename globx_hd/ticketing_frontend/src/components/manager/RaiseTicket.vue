<template>
  <TicketForm @success="handleTicketCreated" />
</template>

<script setup>
import { useRoute } from 'vue-router';
import TicketForm from '../shared/TicketForm.vue';
import { updateDumpedQueryStatus } from '@/api/dumpedQueries';

const route = useRoute();

async function handleTicketCreated(ticket) {
  const dumpId = route.query.dump_id;
  if (dumpId) {
    try {
      console.log(`Resolving dumped query ${dumpId} after ticket creation...`);
      await updateDumpedQueryStatus(dumpId, 'RESOLVED');
      console.log(`Dumped query ${dumpId} resolved successfully.`);
    } catch (e) {
      console.error("Failed to auto-resolve dumped query:", e);
    }
  }
}
</script>
