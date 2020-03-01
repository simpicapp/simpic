import {onMounted, onUnmounted} from "@vue/composition-api";
import Vue from "vue";

export const EventBus = new Vue();

export function useEventListener(event: string, handler: Function) {
  onMounted(() => EventBus.$on(event, handler));
  onUnmounted(() => EventBus.$off(event, handler));
}
