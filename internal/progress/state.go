package progress

import (
	"bytes"

	log "github.com/hashicorp/go-hclog"
	api "github.com/moby/buildkit/api/services/control"
	"github.com/tonistiigi/vt100"
)

type vertexStates struct {
	byTime   []*vertexState
	byDigest map[string]*vertexState
}

type vertexState struct {
	vtx *api.Vertex

	term        *vt100.VT100
	logs        [][]byte
	logsPartial bool

	byID   map[string]*api.VertexStatus
	byTime []*api.VertexStatus
}

func newVertexStates() *vertexStates {
	return &vertexStates{byDigest: make(map[string]*vertexState)}
}

func (s *vertexStates) List() []*vertexState {
	return s.byTime
}

func (s *vertexStates) Get(k string) (*vertexState, bool) {
	v, ok := s.byDigest[k]

	return v, ok
}

func (s *vertexStates) Add(v *api.Vertex, width int) {
	if prev, ok := s.byDigest[v.Digest]; ok {
		if v.Started != nil && prev.vtx.Started == nil {
			s.byTime = append(s.byTime, s.byDigest[v.Digest])
		}

		if prev.vtx.Started == nil || v.Started != nil {
			prev.vtx = v
		}
	} else {
		s.byDigest[v.Digest] = &vertexState{
			vtx:  v,
			byID: make(map[string]*api.VertexStatus),
			term: vt100.NewVT100(termHeight, width),
		}

		if v.Started != nil {
			s.byTime = append(s.byTime, s.byDigest[v.Digest])
		}
	}
}

func (v *vertexState) updateStatus(s *api.VertexStatus) {
	if _, ok := v.byID[s.ID]; !ok {
		v.byTime = append(v.byTime, s)
	}

	v.byID[s.ID] = s
}

func (v *vertexState) updateLogs(l *api.VertexLog) {
	if _, err := v.term.Write(l.Msg); err != nil {
		log.L().Error("error writing log to build terminal", "error", err)
	}

	msg := l.Msg
	isFirst := true

	for {
		if len(msg) == 0 {
			v.logsPartial = false

			return
		}

		index := bytes.IndexByte(msg, byte('\n'))
		if index == -1 {
			v.updateLogLine(msg, isFirst)
			v.logsPartial = true

			return
		}

		v.updateLogLine(msg[:index], isFirst)
		msg = msg[index+1:]
		isFirst = false
	}
}

func (v *vertexState) updateLogLine(line []byte, isFirst bool) {
	if isFirst && v.logsPartial && len(v.logs) != 0 {
		v.logs[len(v.logs)-1] = append(v.logs[len(v.logs)-1], line...)
	} else {
		v.logs = append(v.logs, line)
	}
}
