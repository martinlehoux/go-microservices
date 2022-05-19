import type { Component } from "solid-js";

import Register from "./Register";
import Users from "./Users";

const App: Component = () => {
  return (
    <div class="p-8 flex flex-col items-center">
      <Register />
      <Users />
    </div>
  );
};

export default App;
