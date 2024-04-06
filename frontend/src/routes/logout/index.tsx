import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { removeJWTToken } from "../../service";

export function Logout() {
	removeJWTToken();
	const navigate = useNavigate();

	useEffect(() => {
		navigate("/login");
	}, []);

	return <>Logging Out</>;
}
