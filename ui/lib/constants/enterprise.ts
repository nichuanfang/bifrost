/**
 * Routes backed by enterprise-only UI in the OSS build.
 *
 * Keep this list at the route boundary so enterprise pages cannot be reached
 * by typing a URL or following a stale bookmark. Features that have both OSS
 * and enterprise pieces should hide only their enterprise subcomponents.
 */
const enterpriseOnlyRoutePrefixes = [
	"/workspace/adaptive-routing",
	"/workspace/alerting",
	"/workspace/audit-logs",
	"/workspace/circuit-breaker",
	"/workspace/cluster",
	"/workspace/config/branding",
	"/workspace/config/license",
	"/workspace/config/proxy",
	"/workspace/edge-control",
	"/workspace/governance/access-profiles",
	"/workspace/governance/business-units",
	"/workspace/governance/projects",
	"/workspace/governance/rbac",
	"/workspace/governance/users",
	"/workspace/guardrails",
	"/workspace/mcp-auth-config",
	"/workspace/rbac",
	"/workspace/scim",
] as const;

export function isEnterpriseOnlyRoute(pathname: string): boolean {
	const normalizedPath = pathname.replace(/\/+$/, "") || "/";
	return enterpriseOnlyRoutePrefixes.some((prefix) => normalizedPath === prefix || normalizedPath.startsWith(`${prefix}/`));
}

export { enterpriseOnlyRoutePrefixes };