import {onMounted, onUnmounted} from "@vue/composition-api";

export function useDocumentListener<K extends keyof DocumentEventMap>(type: K, listener: (this: Document, ev: DocumentEventMap[K]) => void) {
  onMounted(() => document.addEventListener(type, listener));
  onUnmounted(() => document.removeEventListener(type, listener));
}

export function useWindowListener<K extends keyof WindowEventMap>(type: K, listener: (this: Window, ev: WindowEventMap[K]) => void) {
  onMounted(() => window.addEventListener(type, listener));
  onUnmounted(() => window.removeEventListener(type, listener));
}
