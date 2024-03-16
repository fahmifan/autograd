import { ConnectError } from "@bufbuild/connect";
import { Anchor, Button, Card, Container, Group, PasswordInput, Stack, } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import { Form, redirect, useNavigate } from "react-router-dom";
import { AutogradServiceClient } from "../../service";

export function AccountActivation() {    
    // get url query params
    const urlParams = new URLSearchParams(window.location.search);
    const activationToken = urlParams.get('activationToken')
    const userID = urlParams.get('userID')

    const [password, setPassword] = useState("")
    const [passwordConfirmation, setPasswordConfirmation] = useState("")
    const navigate = useNavigate();

    async function activateManagedUser(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault()


        try {
            await AutogradServiceClient.activateManagedUser({
                activationToken: activationToken ?? "",
                userId: userID ?? "",
                password,
                passwordConfirmation,
            })
            notifications.show({
                title: "Set Password",
                message: "Success set your password!",
                withCloseButton: true,
                color: "green",
            });
            navigate("/login")
        } catch (err) {
            const err2 = err as ConnectError
            notifications.show({
                title: "Set Password",
                message: `Failed: ${err2.rawMessage ?? "Unknown error"}`,
                withCloseButton: true,
                color: "red",
            });
        }
    }

    return (
        <main>
            <Container maw={400} mt="lg">
                <h1>Set your password</h1>
                <Card shadow="sm" p="lg" radius="sm">

                    <Form onSubmit={activateManagedUser}>
                        <Stack>
                            <PasswordInput
                                required
                                label="Password"
                                placeholder="Your password"
                                radius="md"
                                name="password"
                                id="password"
                                value={password}
                                onChange={(val) => {setPassword(val.target.value)}}
                            />

                            <PasswordInput
                                required
                                label="Confirm Password"
                                placeholder="Your confirm password"
                                radius="md"
                                name="confirm_password"
                                id="confirm_password"
                                value={passwordConfirmation}
                                onChange={(val) => {setPasswordConfirmation(val.target.value)}}
                            />
                        </Stack>

                        <Group justify="space-between" mt="xl">
                            <Anchor />
                            <Button type="submit" radius="xl">
                                Save
                            </Button>
                        </Group>
                    </Form>
                </Card>
            </Container>
        </main>
    );
}