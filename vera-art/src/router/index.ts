import { createRouter, createWebHistory } from "vue-router";
//import HomeView from "../views/HomeView.vue";
//max

import HomeView from "../views/HomeView.vue";
import ArtistView from "../views/whoami/artist.vue";
import BrieflyView from "../views/whoami/briefly.vue";
import IllustratorView from "../views/whoami/illustrator.vue";
import TeacherView from "../views/whoami/teacher.vue";
import VolunteerView from "../views/whoami/volunteer.vue";
import AllWorksView from "../views/menu/all_works.vue";
import NewsView from "../views/menu/news.vue";
import PayDeliveryView from "../views/menu/pay_deliver.vue";
import ServicesView from "../views/menu/services.vue";
import NewsItemView from "../components/News/NewsItemView.vue";
//end_max

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    // {
    //   path: "/about",
    //   name: "about",
    //   component: () => import("../views/AboutView.vue"),
    // },
    {
      path: "/artist",
      name: "artist",
      component: ArtistView,
    },
    {
      path: "/briefly",
      name: "briefly",
      component: BrieflyView,
    },
    {
      path: "/illustrator",
      name: "illustrator",
      component: IllustratorView,
    },
    {
      path: "/teacher",
      name: "teacher",
      component: TeacherView,
    },
    {
      path: "/volunteer",
      name: "volunteer",
      component: VolunteerView,
    },
    {
      path: "/allworks",
      name: "allworks",
      component: AllWorksView,
    },
    {
      path: "/news",
      name: "news",
      component: NewsView,
    },
    {
      path: "/paydelivery",
      name: "paydelivery",
      component: PayDeliveryView,
    },
    {
      path: "/services",
      name: "services",
      component: ServicesView,
    },
    {
      path: "/news/newsitem",
      name: "newsitem",
      component: NewsItemView,
    },
  ],
});

export default router;
