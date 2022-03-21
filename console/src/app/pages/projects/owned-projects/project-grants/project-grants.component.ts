import { animate, state, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { ActivatedRoute, Router } from '@angular/router';
import { tap } from 'rxjs/operators';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { GrantedProject, ProjectGrantState, Role } from 'src/app/proto/generated/zitadel/project_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { ProjectGrantsDataSource } from './project-grants-datasource';

const ROUTEPARAM = 'projectid';

@Component({
  selector: 'cnsl-project-grants',
  templateUrl: './project-grants.component.html',
  styleUrls: ['./project-grants.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*' })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class ProjectGrantsComponent implements OnInit, AfterViewInit {
  public projectId: string = '';
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatTable) public table!: MatTable<GrantedProject.AsObject>;
  public dataSource!: ProjectGrantsDataSource;
  public selection: SelectionModel<GrantedProject.AsObject> = new SelectionModel<GrantedProject.AsObject>(true, []);
  public memberRoleOptions: Role.AsObject[] = [];
  public displayedColumns: string[] = ['grantedOrgName', 'state', 'creationDate', 'changeDate', 'roleNamesList', 'actions'];

  ProjectGrantState: any = ProjectGrantState;

  constructor(
    private mgmtService: ManagementService,
    private toast: ToastService,
    private route: ActivatedRoute,
    private dialog: MatDialog,
    private breadcrumbService: BreadcrumbService,
    private router: Router,
  ) {
    const projectId = this.route.snapshot.paramMap.get(ROUTEPARAM);
    if (projectId) {
      this.projectId = projectId;

      const breadcrumbs = [
        new Breadcrumb({
          type: BreadcrumbType.IAM,
          name: 'IAM',
          routerLink: ['/system'],
        }),
        new Breadcrumb({
          type: BreadcrumbType.ORG,
          routerLink: ['/org'],
        }),
        new Breadcrumb({
          type: BreadcrumbType.PROJECT,
          name: '',
          param: { key: ROUTEPARAM, value: projectId },
          routerLink: ['/projects', projectId],
        }),
      ];
      this.breadcrumbService.setBreadcrumb(breadcrumbs);
    }
  }

  public gotoRouterLink(rL: any) {
    this.router.navigate(rL);
  }

  public ngOnInit(): void {
    this.dataSource = new ProjectGrantsDataSource(this.mgmtService, this.toast);
    this.dataSource.loadGrants(this.projectId, 0, 25, 'asc');
    this.getRoleOptions(this.projectId);
  }

  public ngAfterViewInit(): void {
    this.paginator.page.pipe(tap(() => this.loadGrantsPage())).subscribe();
  }

  public loadGrantsPage(pageIndex?: number, pageSize?: number): void {
    this.dataSource.loadGrants(this.projectId, pageIndex ?? this.paginator.pageIndex, pageSize ?? this.paginator.pageSize);
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.grantsSubject.value.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected()
      ? this.selection.clear()
      : this.dataSource.grantsSubject.value.forEach((row) => this.selection.select(row));
  }

  public getRoleOptions(projectId: string): void {
    this.mgmtService.listProjectRoles(projectId, 100, 0).then((resp) => {
      this.memberRoleOptions = resp.resultList;
    });
  }

  public updateRoles(grant: GrantedProject.AsObject, selectionChange: MatSelectChange): void {
    this.mgmtService
      .updateProjectGrant(grant.grantId, grant.projectId, selectionChange.value)
      .then(() => {
        this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTCHANGED', true);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public deleteGrant(grant: GrantedProject.AsObject): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'PROJECT.GRANT.DIALOG.DELETE_TITLE',
        descriptionKey: 'PROJECT.GRANT.DIALOG.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtService
          .removeProjectGrant(grant.grantId, grant.projectId)
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.REMOVED', true);
            const data = this.dataSource.grantsSubject.getValue();
            this.selection.selected.forEach((item) => {
              const index = data.findIndex((i) => i.grantId === item.grantId);
              if (index > -1) {
                data.splice(index, 1);
                this.dataSource.grantsSubject.next(data);
              }
            });
            this.selection.clear();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }
}
