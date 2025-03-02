package kamailio_cfg

import (
	"KamaiZen/document_manager"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
)

type RouteType = string

const (
	RouteTypeUnknown RouteType = "unknown_route"
	RouteTypeRequest RouteType = "request_route"
	RouteTypeReply   RouteType = "reply_route"
	RouteTypeBranch  RouteType = "branch_route"
	RouteTypeFailure RouteType = "failure_route"
	RouteTypeOnreply RouteType = "onreply_route"
	RouteTypeStartup RouteType = "startup_route"
	RouteTypeNamed   RouteType = "route"
	RouteTypeEvent   RouteType = "event_route"
)

const _REQUEST_ROUTE_DECLARATION_QUERY = `(routing_block
    route: (predef_route) @route
    body: (compound_statement) @body
    (#eq? @route "request_route")
    ) @definition.function`

const _REPLY_ROUTE_DECLARATION_QUERY = `(routing_block
    route: (predef_route) @route
    body: (compound_statement) @body
    (#eq? @route "reply_route")
    ) @definition.function`

const _STARTUP_ROUTE_DECLARATION_QUERY = `(routing_block
    route: (predef_route) @route
    body: (compound_statement) @body
    (#eq? @route "startup_route")
    ) @definition.function`

const _NAMED_ROUTE_DECLARATION_QUERY = `(routing_block
    route: (predef_route) @route
    route_name: (_) @name
    body: (compound_statement) @body
    (#match? @route "route|branch_route|failure_route|onreply_route|onsend_route|event_route")
    ) @definition.function`

type Route struct {
	doc        string
	Content    string
	StartPoint sitter.Point
	EndPoint   sitter.Point
	routeType  RouteType
	Name       string
}

type Routes struct {
	requestRoute  Route
	replyRoute    Route
	startupRoute  Route
	branchRoutes  []Route
	failureRoutes []Route
	onreplyRoutes []Route
	onSendRoutes  []Route
	eventRoutes   []Route
	routes        []Route
}

func (r *Routes) ContainsRequestRoute() bool {
	return r.requestRoute.Content != ""
}

func (r *Routes) AddNamedRoute(route ...Route) {
	for _, rt := range route {
		switch rt.routeType {
		case RouteTypeBranch:
			r.branchRoutes = append(r.branchRoutes, rt)
		case RouteTypeFailure:
			r.failureRoutes = append(r.failureRoutes, rt)
		case RouteTypeOnreply:
			r.onreplyRoutes = append(r.onreplyRoutes, rt)
		case RouteTypeEvent:
			r.eventRoutes = append(r.eventRoutes, rt)
		case RouteTypeNamed:
			r.routes = append(r.routes, rt)
		}
	}
}

func (r *Routes) FindRouteByNameAndType(name string, routeType RouteType) *Route {
	getRouteByName := func(routes []Route, rn string) *Route {
		for _, route := range routes {
			if route.Name == rn {
				return &route
			}
		}
		return nil
	}
	switch routeType {
	case RouteTypeBranch:
		return getRouteByName(r.branchRoutes, name)
	case RouteTypeFailure:
		return getRouteByName(r.failureRoutes, name)
	case RouteTypeOnreply:
		return getRouteByName(r.onreplyRoutes, name)
	case RouteTypeEvent:
		return getRouteByName(r.eventRoutes, name)
	case RouteTypeNamed:
		return getRouteByName(r.routes, name)
	default:
		return nil
	}
}

func FetchRoutes(a *Analyzer, source_code []byte) Routes {
	var routes Routes
	r := fetchRequestRoute(a, source_code)
	if r != nil {
		routes.requestRoute = *r
	}
	r = fetchRequestRoute(a, source_code)
	if r != nil {
		routes.replyRoute = *r
	}
	r = fetchRequestRoute(a, source_code)
	if r != nil {
		routes.requestRoute = *r
	}
	r = fetchReplyRoute(a, source_code)
	if r != nil {
		routes.replyRoute = *r
	}
	routes.AddNamedRoute(fetchNamedRoutes(a, source_code)...)
	return routes
}

func fetchRequestRoute(a *Analyzer, source_code []byte) *Route {
	q, err := NewQueryExecutor(
		_REQUEST_ROUTE_DECLARATION_QUERY,
		a.ast.Node,
		a.builder.parser.language,
	)
	if err != nil {
		return nil
	}

	_routeTag := "route"
	_bodyTag := "body"

	var route *Route
	var routeType RouteType
	var routeBody string
	var start, end sitter.Point
	match, ok := q.NextMatch()
	if !ok {
		// TODO: think
		// This is not fine since we expect at least one request route
		// however only one file should have it
		log.Error().Msg("No request route found")
	} else {
		for _, capture := range match.Captures {
			node := capture.Node
			captureName := q.query.CaptureNameForId(capture.Index)
			if captureName == _routeTag {
				routeType = RouteTypeRequest
				start = node.StartPoint()
			}
			if captureName == _bodyTag {
				routeBody = string(node.Content(source_code))
				end = node.EndPoint()
			}

		}
		route = &Route{
			doc:        document_manager.GetCookBookDocs("request_route"),
			Name:       "",
			Content:    routeBody,
			StartPoint: start,
			EndPoint:   end,
			routeType:  routeType,
		}
		match, ok = q.NextMatch()
		if ok {
			// TODO:
			// File should have only one request route
			// may be note the position and add a diagnostic
			log.Error().Msg("More than one request route found")
		}
	}
	return route
}

func fetchReplyRoute(a *Analyzer, source_code []byte) *Route {
	q, err := NewQueryExecutor(
		_REPLY_ROUTE_DECLARATION_QUERY,
		a.ast.Node,
		a.builder.parser.language,
	)
	if err != nil {
		return nil
	}

	_routeTag := "route"
	_bodyTag := "body"

	var routeType RouteType
	var routeBody string
	var start, end sitter.Point
	match, ok := q.NextMatch()
	if !ok {
		return nil
	}
	for _, capture := range match.Captures {
		node := capture.Node
		captureName := q.query.CaptureNameForId(capture.Index)
		if captureName == _routeTag {
			routeType = RouteTypeRequest
			start = node.StartPoint()
		}
		if captureName == _bodyTag {
			routeBody = string(node.Content(source_code))
			end = node.EndPoint()
		}
	}
	return &Route{
		Name:       "",
		Content:    routeBody,
		StartPoint: start,
		EndPoint:   end,
		routeType:  routeType,
	}
}

func fetchNamedRoutes(a *Analyzer, source_code []byte) []Route {
	var routes []Route

	q, err := NewQueryExecutor(
		_NAMED_ROUTE_DECLARATION_QUERY,
		a.ast.Node,
		a.builder.parser.language,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error creating query executor for named route")
		return nil
	}
	_routeTag := "route"
	_bodyTag := "body"
	_nameTag := "name"

	var routeName string
	var routeType RouteType
	var routeBody string
	var start, end sitter.Point
	for {
		match, ok := q.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			captureName := q.query.CaptureNameForId(capture.Index)
			if captureName == _routeTag {
				routeType = capture.Node.Content(source_code)
				start = node.StartPoint()
			}

			if captureName == _nameTag {
				routeName = capture.Node.Content(source_code)
			}

			if captureName == _bodyTag {
				routeBody = string(node.Content(source_code))
				end = node.EndPoint()
			}

		}
		// TODO: use the doxygen docs if available
		route := Route{
			doc:        document_manager.GetCookBookDocs(routeType),
			Name:       routeName,
			Content:    routeBody,
			StartPoint: start,
			EndPoint:   end,
			routeType:  routeType,
		}
		routes = append(routes, route)

	}
	return routes

}
