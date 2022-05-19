import { users } from "./resources/users";

export default function Users() {
  return (
    <main>
      <h1 class="font-bold text-xl">Users</h1>
      <table class="table-fixed w-full">
        <thead>
          <tr>
            <th>Email</th>
            <th>Preferred name</th>
          </tr>
        </thead>
        <tbody>
          {users()?.map((user) => (
            <tr>
              <td>{user.email}</td>
              <td>{user.preferred_name}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </main>
  );
}
