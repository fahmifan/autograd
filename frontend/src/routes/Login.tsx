import { ActionFunctionArgs, Form, redirect } from "react-router-dom";
import { AutogradServiceClient } from "../service";

export default function Login() {
    return (
        <div>
            <h1>Login</h1>
            <Form method="post" id="login-form">
                <p>
                    <label htmlFor="email">Email</label>
                    <input type="email" name="email" id="email" />
                </p>
                <p>
                    <label htmlFor="password">Password</label>
                    <input type="password" name="password" id="password" />
                </p>

                <button type="submit">Login</button>
            </Form>
        </div>
    );
}

export async function loginAction({ request }: ActionFunctionArgs): Promise<Response | null> {
    const formData = await request.formData()
    const email = formData.get("email") as string
    const password = formData.get("password") as string

    const res = await AutogradServiceClient.login({
        email,
        password
    })

    if (res) {
        localStorage.setItem("token", res.token)
        return redirect("/user-management")
    }

    return null;
}