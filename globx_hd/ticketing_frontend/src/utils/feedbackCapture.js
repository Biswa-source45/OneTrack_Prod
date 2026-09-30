import { ref } from 'vue';
import { toBlob } from 'html-to-image';

// Hand-off from the Ctrl+I capture (Layout.vue) to the Feedback page. A route
// change can't carry a Blob, so it waits here until the page consumes it.
export const feedbackDraft = ref(null); // { file, title, pageUrl }

export const CAPTURE_OVERLAY_ID = 'feedback-capture-overlay';

export async function capturePageScreenshot() {
  const blob = await toBlob(document.body, {
    type: 'image/jpeg',
    quality: 0.85,
    pixelRatio: 1,            // keeps full-page captures well under 5 MB
    backgroundColor: '#ffffff', // JPEG has no transparency; avoid black areas
    filter: (node) => node.id !== CAPTURE_OVERLAY_ID,
  });
  if (!blob) throw new Error('capture failed');
  return new File([blob], `screenshot-${Date.now()}.jpg`, { type: 'image/jpeg' });
}

// "/manager/master-data/products/12/edit" -> "Products"
export function pageNameFromPath(path) {
  const parts = path.split('/').filter(p => p && !/^\d+$/.test(p) && !['edit', 'new', 'create'].includes(p));
  if (parts.length > 1 && ['manager', 'engineer', 'contacts'].includes(parts[0])) parts.shift();
  const part = parts.pop();
  if (!part) return 'Dashboard';
  return part.split('-').map(w => w[0].toUpperCase() + w.slice(1)).join(' ');
}
