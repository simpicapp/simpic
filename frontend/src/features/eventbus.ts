import {onMounted, onUnmounted} from "@vue/composition-api";
import {EventBus} from "@/components/bus";

export function useEventListener(event: string, handler: Function) {
  onMounted(() => EventBus.$on(event, handler));
  onUnmounted(() => EventBus.$off(event, handler));
}
