import {computed, ref} from "@vue/composition-api";
import {AxiosError} from "axios";

export function useAlert() {
  const alert = ref("");
  const hasAlert = computed(() => alert.value.length > 0);

  function setAlert(error?: AxiosError) {
    if (!error) {
      alert.value = "";
    } else if (error.response) {
      alert.value = error.response.data.error || error.message;
    } else {
      alert.value = error.message;
    }
  }

  return {alert, hasAlert, setAlert};
}
