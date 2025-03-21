/* @refresh reload */
import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";

import "./index.css";
import Home from "./pages/home";
import Layout from "./layout";

const root = document.getElementById("root");

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    "Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?",
  );
}

render(
  () => (
    <Router root={Layout}>
      <Route path="/" component={Home} />
    </Router>
  ),
  root!,
);
