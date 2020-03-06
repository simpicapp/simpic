import VueRouter from "vue-router";
import Lightbox from "@/pages/lightbox.vue";
import Timeline from "@/pages/timeline.vue";
import Albums from "@/pages/albums.vue";
import Album from "@/pages/album.vue";

const router = new VueRouter({
  mode: "history",
  routes: [
    {
      children: [
        {
          component: Lightbox,
          path: "/timeline/photo/:id",
          props: true,
        },
      ],
      component: Timeline,
      path: "/timeline/",
    },
    {
      component: Albums,
      path: "/albums/",
    },
    {
      children: [
        {
          component: Lightbox,
          path: "/albums/:album/photo/:id",
          props: true,
        },
      ],
      component: Album,
      path: "/albums/:id",
      props: true,
    },
    {
      path: "/",
      redirect: "/timeline/",
    },
  ],
});

router.beforeEach((to, _, next) => {
  if (to.fullPath.endsWith("/")) {
    next();
  } else {
    next(to.fullPath + "/");
  }
});

export function useRouter() {
  return {router};
}
