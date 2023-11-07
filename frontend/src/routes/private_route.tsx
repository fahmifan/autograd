import { Navigate, useLocation } from "react-router-dom";
import { getDecodedJWTToken } from "../service";

export function PrivateRoute({ element }: { element: JSX.Element }) {
	const jwtToken = getDecodedJWTToken();
	const location = useLocation();

	if (!jwtToken || !jwtToken.id) {
		return <Navigate to="/login" state={{ from: location }} />
	}


	return element
}
