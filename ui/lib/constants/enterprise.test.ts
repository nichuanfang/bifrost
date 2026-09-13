import { describe, expect, it } from "vitest";
import { isEnterpriseOnlyRoute } from "./enterprise";

describe("isEnterpriseOnlyRoute", () => {
	it("matches enterprise routes and nested pages", () => {
		expect(isEnterpriseOnlyRoute("/workspace/cluster")).toBe(true);
		expect(isEnterpriseOnlyRoute("/workspace/cluster/nodes")).toBe(true);
		expect(isEnterpriseOnlyRoute("/workspace/governance/projects")).toBe(true);
	});

	it("does not match similarly prefixed OSS routes", () => {
		expect(isEnterpriseOnlyRoute("/workspace/clustered-models")).toBe(false);
		expect(isEnterpriseOnlyRoute("/workspace/governance/teams")).toBe(false);
		expect(isEnterpriseOnlyRoute("/workspace/providers")).toBe(false);
	});

	it("normalizes trailing slashes", () => {
		expect(isEnterpriseOnlyRoute("/workspace/scim/")).toBe(true);
	});
});