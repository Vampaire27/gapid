// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/gapid/test/robot/build"
	q "github.com/google/gapid/test/robot/search/query"
)

func (s *Server) handleArtifacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result := []*build.Artifact(nil)

	if query, err := query(w, r); err == nil {
		if err = s.Build.SearchArtifacts(ctx, query, func(ctx context.Context, entry *build.Artifact) error {
			result = append(result, entry)
			return nil
		}); err != nil {
			writeError(w, 500, err)
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

func (s *Server) handlePackages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result := []*build.Package(nil)

	if query, err := query(w, r); err == nil {
		if err = s.Build.SearchPackages(ctx, query, func(ctx context.Context, entry *build.Package) error {
			result = append(result, entry)
			return nil
		}); err != nil {
			writeError(w, 500, err)
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

func (s *Server) handlePackageChain(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	head := r.FormValue("head")
	if head == "" {
		writeError(w, 404, errors.New("The head parameter is required"))
		return
	}

	pkgs := map[string]*build.Package{}
	if err := s.Build.SearchPackages(ctx, q.Bool(true).Query(), func(ctx context.Context, entry *build.Package) error {
		pkgs[entry.Id] = entry
		return nil
	}); err != nil {
		writeError(w, 500, err)
		return
	}

	pkg, ok := pkgs[head]
	if !ok {
		writeError(w, 404, fmt.Errorf("Package '%s' not found", head))
		return
	}

	result := []*build.Package(nil)
	for ; pkg != nil; pkg = pkgs[pkg.Parent] {
		result = append(result, pkg)
	}

	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleTracks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result := []*build.Track(nil)

	if query, err := query(w, r); err == nil {
		if err = s.Build.SearchTracks(ctx, query, func(ctx context.Context, entry *build.Track) error {
			result = append(result, entry)
			return nil
		}); err != nil {
			writeError(w, 500, err)
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}
