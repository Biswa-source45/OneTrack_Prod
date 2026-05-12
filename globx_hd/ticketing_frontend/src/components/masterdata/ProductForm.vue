<template>
  <div class="p-8">
    <PageHeader :title="isEdit ? 'Edit Product' : 'Add Product'" />
    <FormLayout @submit="onSubmit">
      <FormField label="Product Name" :error="errors.product_name">
        <input v-model="form.product_name" type="text" class="w-full border border-blue-200 rounded px-3 py-2" required />
      </FormField>
      <FormField label="Description">
        <input v-model="form.product_description" type="text" class="w-full border border-blue-200 rounded px-3 py-2" />
      </FormField>
      <template #actions>
        <Button variant="secondary" type="button" @click="cancel">Cancel</Button>
        <Button type="submit">Save</Button>
      </template>
    </FormLayout>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import PageHeader from '../ui/PageHeader.vue';
import FormLayout from '../ui/FormLayout.vue';
import FormField from '../ui/FormField.vue';
import Button from '../ui/Button.vue';
import { createProduct, updateProduct, fetchProducts } from '../../api/auth';

const route = useRoute();
const router = useRouter();
const isEdit = route.params.id !== undefined;
const form = reactive({ product_name: '', product_description: '' });
const errors = reactive({ product_name: '' });

onMounted(async () => {
  if (isEdit) {
    const id = route.params.id;
    // Since we don't have GET /products/:id, fetch all and find
    const list = await fetchProducts();
    const item = list.find(p => String(p.id) === String(id));
    if (item) {
      form.product_name = item.product_name || '';
      form.product_description = item.product_description || '';
    }
  }
});

const submitting = ref(false);
async function onSubmit() {
  if (submitting.value) return;
  errors.product_name = form.product_name ? '' : 'Required';
  if (errors.product_name) return;
  submitting.value = true;
  if (isEdit) {
    await updateProduct(route.params.id, { product_name: form.product_name, product_description: form.product_description });
  } else {
    await createProduct({ product_name: form.product_name, product_description: form.product_description });
  }
  router.push('/master-data/products');
}
function cancel() { router.back(); }
</script>
