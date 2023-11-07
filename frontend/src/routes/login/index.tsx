import {
	Anchor,
	Button,
	Container,
	Group,
	PasswordInput,
	Stack,
	Text,
	TextInput,
} from "@mantine/core";
import { ActionFunctionArgs, Form, Navigate, redirect } from "react-router-dom";
import { AutogradServiceClient, getDecodedJWTToken, saveJWTToken } from "../../service";

export function LoginPage() {
	const decoded = getDecodedJWTToken();
	if (decoded?.id) {
		if (decoded?.role === "admin") {
			return <Navigate to="/backoffice" />;
		} else if (decoded?.role === "student") {
			return <Navigate to="/student-dashboard" />;
		}
	}

	return (
		<Container maw={400} pt="md">
			<Text size="lg" fw={500} mb="md">
				Welcome to Autograd, login to continue
			</Text>

			<Form method="POST" id="login-form">
				<Stack>
					<TextInput
						required
						label="Email"
						placeholder="your@email.com"
						radius="md"
						name="email"
						id="email"
					/>

					<PasswordInput
						required
						label="Password"
						placeholder="Your password"
						radius="md"
						name="password"
						id="password"
					/>
				</Stack>

				<Group justify="space-between" mt="xl">
					<Anchor />
					<Button type="submit" radius="xl">
						Login
					</Button>
				</Group>
			</Form>
		</Container>
	);
}

export async function loginAction({
	request,
}: ActionFunctionArgs): Promise<Response | null> {
	const formData = await request.formData();
	const email = formData.get("email") as string;
	const password = formData.get("password") as string;

	const res = await AutogradServiceClient.login({
		email,
		password,
	});

	if (!res) {
		return null;
	}

	saveJWTToken(res.token);
	const decoded = getDecodedJWTToken();

	if (decoded?.role === "admin") {
		return redirect("/backoffice");
	}

	return redirect("/student-dashboard");
}
