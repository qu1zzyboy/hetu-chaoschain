import { request } from '@umijs/max';
import { BASE_URL } from '@/constants/config';
console.log(BASE_URL)
export async function proposals(
  body: AGENTAPI.ProposalsReq,
  options?: { [key: string]: any },
) {
  return request<AGENTAPI.ProposalsRes>('/api/proposals', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function proposalDetail(
  body: AGENTAPI.ProposalDetailReq,
  options?: { [key: string]: any },
) {
  return request<AGENTAPI.ProposalDetailRes>('/api/proposal-detail', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}


export async function manifesto(
  options?: { [key: string]: any },
) {
  return request<AGENTAPI.ManifestoRes>('/api/manifesto', {
    method: 'GET',
    ...(options || {}),
  });
}

export async function agents(
  options?: { [key: string]: any },
) {
  return request<AGENTAPI.AgentsRes>(`${BASE_URL ?? '/api'}/agents`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}

export async function networkStatus(
  options?: { [key: string]: any },
) {
  return request<AGENTAPI.NetworkStatusRes>('/api/network-status', {
    method: 'GET',
    ...(options || {}),
  });
}


export async function latestBlocks(
  options?: { [key: string]: any },
) {
  return request<AGENTAPI.LatestBlocksRes>('/api/latest-blocks', {
    method: 'GET',
    ...(options || {}),
  });
}


export async function agentDetail(
  body: AGENTAPI.AgentDetailReq,
  options?: { [key: string]: any },
) {
  return request<AGENTAPI.AgentDetailRes>('/api/agent-detail', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
