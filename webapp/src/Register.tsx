import { authenticationGateway, humanResourcesGateway } from "./module";
import { mutate, refetch } from "./resources/users";

async function handleSubmit(event: Event) {
  event.preventDefault();

  const form = event.target as HTMLFormElement;
  const identifier = form.elements.namedItem("identifier") as HTMLInputElement;
  const password = form.elements.namedItem("password") as HTMLInputElement;
  const preferred_name = form.elements.namedItem(
    "preferred_name"
  ) as HTMLInputElement;

  await authenticationGateway.register({
    identifier: identifier.value,
    password: password.value,
  });

  const { token } = await authenticationGateway.authenticate({
    identifier: identifier.value,
    password: password.value,
  });
  humanResourcesGateway.setToken(token);

  await humanResourcesGateway.register({
    preferred_name: preferred_name.value,
  });

  form.reset();

  await refetch();
}

export default function Register() {
  return (
    <form
      onSubmit={handleSubmit}
      class="max-w-md flex flex-col shadow rounded p-4 gap-2"
    >
      <input type="text" name="preferred_name" placeholder="Preferred name" />
      <input type="email" name="identifier" placeholder="Email" />
      <input type="password" name="password" placeholder="Password" />
      <button type="submit">Register</button>
    </form>
  );
}
