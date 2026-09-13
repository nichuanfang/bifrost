import { IS_ENTERPRISE } from "@/lib/constants/config";
import { isEnterpriseOnlyRoute } from "@/lib/constants/enterprise";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { ClientLayout } from "../clientLayout";

function WorkspaceLayout({ children }: { children: React.ReactNode }) {
	return <ClientLayout>{children}</ClientLayout>;
}

function RouteComponent() {
	return (
		<WorkspaceLayout>
			<Outlet />
		</WorkspaceLayout>
	);
}

export const Route = createFileRoute("/workspace")({
	beforeLoad: ({ location }) => {
		if (location.pathname === "/workspace" || location.pathname === "/workspace/") {
			throw redirect({ to: "/workspace/dashboard", replace: true });
		}
		if (!IS_ENTERPRISE && isEnterpriseOnlyRoute(location.pathname)) {
			throw redirect({ to: "/workspace/dashboard", replace: true });
		}
	},
	component: RouteComponent,
});