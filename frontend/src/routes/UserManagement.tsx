import { ActionFunctionArgs, Form, Outlet, redirect, useLoaderData } from "react-router-dom";
import { ManagedUser } from "../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../service";

export function ListManagedUsers() {
    const { managedUsers } =  useLoaderData() as LoaderResponse;

    if (!managedUsers || managedUsers.length === 0) {
        return <>
            <h2>List Users</h2>
            <p><i>No Users</i></p>
        </>
    }

    return <>
        <h2>List Users</h2>
        {
            managedUsers.map((user) => {
                return (
                    <div key={user.id}>
                        <p>{user.name}</p>
                        <p>{user.email}</p>
                    </div>
                )
            })
        }
    </>
}

export default function UserManagement() {
  return (
    <div>
        <h1>User Management</h1>
        <Form action="create">
            <button type="submit">Create</button>
        </Form>

        <Outlet />
    </div>
  );
}

type LoaderResponse = {
    managedUsers: ManagedUser[];
}

export async function loaderUserManagement(): Promise<LoaderResponse> {
    const res = await AutogradServiceClient.findAllManagedUsers({
        limit: 10,
        page: 1,
    })

    return res
}

export function CreateManagedUser() {
    return <>
        <h1>User Management</h1>
        <h2>Create User</h2>
        <section>
            <Form method="post" id="create-managed-user">
                <p>
                    <label htmlFor="name">Name</label>
                    <input type="text" name="name" id="name" />
                </p>
                <p>
                    <label htmlFor="email">Email</label>
                    <input type="email" name="email" id="email" />
                </p>
                <p>
                    <label htmlFor="role">Role</label>
                    <select name="role" id="role" placeholder="Select role">
                        <option value="admin">Admin</option>
                        <option value="teacher">Teacher</option>
                        <option value="student">Student</option>
                    </select>
                </p>

                <button type="submit">Create User</button>
            </Form>
        </section>
    </>
}

export async function actionCreateManagedUser({ request }: ActionFunctionArgs): Promise<Response | null> {
    const formData = await request.formData()
    
    const name = formData.get("name") as string
    const email = formData.get("email") as string
    const role = formData.get("role") as string

    const res = await AutogradServiceClient.createManagedUser({
        email,
        name,
        role,
    })

    if (res) {
        return redirect("/user-management")
    }

    return null;
}
